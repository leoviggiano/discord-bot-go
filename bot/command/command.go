package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/bot/command/ping"
	"bot_test/entity"
	"bot_test/lib/discord"
	"bot_test/service"
)

type Command interface {
	Name() string
	Description() string
	Options() []*discordgo.ApplicationCommandOption
	Cooldown() int

	ExecEvent(i *discordgo.InteractionCreate, customID *discord.CustomID) (entity.Payload, error)
}

func All(s *discordgo.Session, services service.All, log logrus.FieldLogger) []Command {
	return []Command{
		ping.NewCommand(s, services, log),
	}
}
