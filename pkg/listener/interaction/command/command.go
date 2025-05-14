package command

import (
	"github.com/bwmarrin/discordgo"

	"megumin/pkg/listener/interaction/command/ping"
	"megumin/pkg/listener/interaction/context"
	"megumin/pkg/service"
)

type Command interface {
	Name() string
	Description() string
	Permission() int
	Cooldown() int
	Options() []*discordgo.ApplicationCommandOption

	Events() map[string]func(*context.Context) (string, error)
}

func All(services service.All) []Command {
	return []Command{
		ping.NewCommand(services),
	}
}
