package ping

import (
	"github.com/bwmarrin/discordgo"

	"megumin/pkg/constants/events"
	"megumin/pkg/listener/interaction/context"
	"megumin/pkg/service"
)

type command struct {
	services service.All
	events   map[string]func(*context.Context) (string, error)
}

type commandPayload struct {
}

func NewCommand(services service.All) command {
	c := command{
		services: services,
	}

	c.events = map[string]func(*context.Context) (string, error){
		c.Name(): c.InteractionCreate,
	}

	return c
}

func (c command) Name() string {
	return events.Ping
}

func (c command) Description() string {
	return "Ping the bot"
}

func (c command) Permission() int {
	return 0
}

func (c command) Cooldown() int {
	return 0
}

func (c command) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c command) Events() map[string]func(*context.Context) (string, error) {
	return c.events
}

func (c command) Payload(payload any) *commandPayload {
	p, ok := payload.(*commandPayload)
	if !ok {
		p = &commandPayload{}
	}

	return p
}
