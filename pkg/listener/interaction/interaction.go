package interaction

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"megumin/config"
	"megumin/pkg/bot/dserr"
	"megumin/pkg/constants/permissions"
	"megumin/pkg/lib/discord"
	"megumin/pkg/lib/sql"
	"megumin/pkg/listener/interaction/command"
	"megumin/pkg/listener/interaction/context"
	"megumin/pkg/listener/interaction/middleware"
	"megumin/pkg/service"
)

type interactionCreate struct {
	log         *slog.Logger
	commands    map[string]command.Command
	events      map[string]func(*context.Context) (string, error)
	middlewares []middleware.Middleware

	services service.All
	session  discord.Session

	onError func(*discordgo.InteractionCreate, error)
}

func NewListener(s *discordgo.Session, db sql.Connection, log *slog.Logger) (*interactionCreate, error) {
	services := service.GetAll(db, log)
	session := discord.NewSession(s)
	errorHandler := dserr.NewErrorHandler(session, log)

	interaction := &interactionCreate{
		commands:    make(map[string]command.Command),
		events:      make(map[string]func(*context.Context) (string, error)),
		middlewares: middleware.All(),
		services:    services,
		log:         log,
		session:     session,

		onError: func(i *discordgo.InteractionCreate, err error) {
			embedError := errorHandler.Error(i, err)

			if errors.Is(err, discord.ErrWithFollowUp) {
				session.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: embedError.Data.Embeds,
					Flags:  embedError.Data.Flags,
				})

				return
			}
			interactionErr := session.InteractionRespond(i.Interaction, embedError)
			if interactionErr != nil {
				log.Error("Error responding to interaction", "error", interactionErr.Error())
			}
		},
	}

	err := interaction.registerCommands(session, services)
	if err != nil {
		return nil, err
	}

	return interaction, nil
}

func (l *interactionCreate) registerCommands(session discord.Session, services service.All) error {
	registeredCommands, err := session.ApplicationCommands(config.AppID(), config.GuildID())
	if err != nil {
		return err
	}

	mappedCommands := make(map[string]bool)

	for _, cmd := range command.All(services) {
		mappedCommands[cmd.Name()] = true
		l.registerCommandEvents(cmd.Events())

		l.commands[cmd.Name()] = cmd
		c := &discordgo.ApplicationCommand{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Type:        discordgo.ChatApplicationCommand,
			Options:     cmd.Options(),
		}

		if !hasCommand(registeredCommands, cmd.Name()) {
			guildID := config.GuildID()
			if cmd.Permission() == permissions.Admin {
				guildID = config.AdminGuildID()
			}

			_, err := session.ApplicationCommandCreate(config.AppID(), guildID, c)
			if err != nil {
				fmt.Println(err)
				continue
			}

			l.log.Info("Command registered with success", "command", cmd.Name())
		}
	}

	// for _, cmd := range registeredCommands {
	// 	if _, ok := mappedCommands[cmd.Name]; !ok {
	// 		err := session.ApplicationCommandDelete(config.AppID(), config.GuildID(), cmd.ID)
	// 		if err != nil {
	// 			l.log.Error("Error deleting command", "error", err.Error())
	// 		}
	// 		l.log.Info("Command removed with success", "command", cmd.Name)
	// 	}
	// }

	// _, err = session.ApplicationCommandBulkOverwrite(config.AppID(), config.GuildID(), registeredCommands)
	// if err != nil {
	// 	l.log.Error(err)
	// }

	return nil
}

func (i interactionCreate) registerCommandEvents(events map[string]func(*context.Context) (string, error)) {
	for key, event := range events {
		if _, ok := i.events[key]; ok {
			fmt.Printf("event %s already registered\n", key)
			continue
		}

		i.events[key] = event
	}
}

func hasCommand(listCommands []*discordgo.ApplicationCommand, commandName string) bool {
	for _, cmd := range listCommands {
		if cmd.Name == commandName {
			return true
		}
	}

	return false
}
