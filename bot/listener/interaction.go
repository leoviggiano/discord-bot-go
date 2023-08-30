package listener

import (
	"fmt"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/bot/command"
	"bot_test/bot/dserr"
	"bot_test/bot/middleware"
	"bot_test/config"
	"bot_test/constants/errors"
	"bot_test/entity"
	"bot_test/lib/discord"
	"bot_test/service"
)

type Interaction interface {
	Handler(*discordgo.Session, *discordgo.InteractionCreate)
}

type interaction struct {
	log         logrus.FieldLogger
	commands    map[string]command.Command
	middlewares []middleware.Middleware
	OnError     func(*discordgo.InteractionCreate, error)
	services    service.All
}

func NewInteraction(session *discordgo.Session, services service.All, log logrus.FieldLogger) Interaction {
	interaction := &interaction{
		commands:    make(map[string]command.Command),
		middlewares: middleware.All(),
		services:    services,
		log:         log,
		OnError: func(i *discordgo.InteractionCreate, err error) {
			interactionErr := session.InteractionRespond(i.Interaction, dserr.Error(session, i, services, log, err))
			if interactionErr != nil {
				log.Error(interactionErr)
			}
		},
	}

	interaction.registerCommands(session, services)
	return interaction
}

func (l interaction) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if panic := recover(); panic != nil {
			fmt.Println(string(debug.Stack()))
			l.OnError(i, panic.(error))
		}
	}()

	var err error
	customID := discord.CustomIDFromInteraction(i)

	if i.Type == discordgo.InteractionApplicationCommand {
		commandName := i.ApplicationCommandData().Name

		customID, err = l.NewCustomID(commandName)
		if err != nil {
			l.OnError(i, err)
			return
		}
	}

	cmd := l.commands[customID.CommandName]
	if cmd == nil {
		l.OnError(i, errors.ErrCommandNotFound)
		return
	}

	payload, err := cmd.ExecEvent(i, customID)
	if err != nil {
		fmt.Println(err)
		l.OnError(i, err)
		return
	}

	err = l.services.CustomID.Update(payload)
	if err != nil {
		l.OnError(i, err)
		return
	}
}

func (l *interaction) registerCommands(session *discordgo.Session, services service.All) {
	commands := make([]*discordgo.ApplicationCommand, 0)

	registeredCommands, err := session.ApplicationCommands(config.AppID(), config.GuildID())
	if err != nil {
		fmt.Println(err)
	}

	for _, cmd := range command.All(session, services, l.log) {
		l.commands[cmd.Name()] = cmd
		c := &discordgo.ApplicationCommand{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Options:     cmd.Options(),
			Type:        discordgo.ChatApplicationCommand,
		}
		commands = append(commands, c)

		if !hasCommand(registeredCommands, cmd.Name()) {
			guildID := config.GuildID()
			_, err := session.ApplicationCommandCreate(config.AppID(), guildID, c)
			if err != nil {
				fmt.Println(err)
				continue
			}

			l.log.Infof("%s command registered with success", cmd.Name())
		}
	}

	// for _, cmd := range registeredCommands {
	// 	fmt.Println(cmd)
	// 	err := session.ApplicationCommandDelete(config.AppID(), config.GuildID(), cmd.ID)
	// 	fmt.Println(err)
	// }

	// _, err = session.ApplicationCommandBulkOverwrite(config.AppID(), config.GuildID(), commands)
	// if err != nil {
	// 	l.log.Error(err)
	// }
}

func hasCommand(listCommands []*discordgo.ApplicationCommand, commandName string) bool {
	for _, cmd := range listCommands {
		if cmd.Name == commandName {
			return true
		}
	}

	return false
}

func (l *interaction) NewCustomID(commandName string) (*discord.CustomID, error) {
	payload := &entity.BasePayload{
		CustomID: discord.NewCustomID(commandName),
	}

	_, err := l.services.CustomID.Set(payload)
	return payload.CustomID, err
}
