package interaction

import (
	"fmt"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/constants/errors"
	"megumin/pkg/constants/events"
	"megumin/pkg/lib/discord"
	"megumin/pkg/listener/interaction/context"
	"megumin/pkg/listener/interaction/middleware"
)

func (l interactionCreate) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := discord.UserFromInteraction(i)
	defer func() {
		if panic := recover(); panic != nil {
			l.log.Error(string(debug.Stack()))
			l.onError(i, panic.(error))
		}
	}()

	customID := l.GetCustomID(i)
	payload, err := l.services.CustomID.Get(customID)
	if err != nil {
		l.onError(i, err)
		return
	}

	ephemeral := i.Message != nil && i.Message.Flags&discordgo.MessageFlagsEphemeral != 0
	if !payload.React(user.ID, ephemeral) {
		return
	}

	if payload.Event == events.CommandEnded {
		l.session.InteractionRespond(i.Interaction, discord.InteractionNoResponse())
		return
	}

	cmd := l.commands[payload.Command]
	if cmd == nil {
		l.onError(i, errors.ErrCommandNotFound)
		return
	}

	ctx := &context.Context{
		User:              user,
		InteractionCreate: i,
		Payload:           payload,
		CustomID:          customID,
		Session:           discord.NewSession(s),
		Log:               l.log,
	}

	for _, m := range l.middlewares {
		next, err := m.Exec(middleware.Args{
			Command:           cmd,
			Session:           discord.NewSession(s),
			InteractionCreate: i,
			Services:          l.services,
			Event:             payload.Event,
			CustomID:          customID,
		})

		if err != nil {
			l.onError(i, err)
			return
		}

		if !next {
			return
		}
	}

	eventName := cmd.Name()
	if payload.Event != "" {
		eventName = payload.Event
	}

	execEvent, ok := l.events[eventName]
	if !ok {
		l.onError(i, fmt.Errorf("%w: %s - CustomID: %s", errors.ErrEventNotFound, eventName, customID))
		return
	}

	nextEvent, err := execEvent(ctx)
	if err != nil {
		l.onError(i, err)
		return
	}

	if nextEvent != "" {
		payload.Event = nextEvent
	}
}

func (l *interactionCreate) GetCustomID(i *discordgo.InteractionCreate) discord.CustomID {
	customID := discord.PayloadKeyFromInteraction(i)

	if i.Type == discordgo.InteractionApplicationCommand {
		commandName := i.ApplicationCommandData().Name
		payload := &context.Payload{
			Command:  commandName,
			Reacters: make([]string, 0),
		}

		return l.services.CustomID.Set(payload)
	}

	return discord.CustomID(customID)
}
