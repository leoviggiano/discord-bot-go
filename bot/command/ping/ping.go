package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/constants/errors"
	"bot_test/entity"
	"bot_test/lib/discord"
	"bot_test/service"
)

const (
	eventConfirm = "confirm"
	eventCancel  = "cancel"
)

type command struct {
	session  *discordgo.Session
	services service.All
	log      logrus.FieldLogger
	events   map[string]func(*discordgo.InteractionCreate, *payload) (entity.Payload, error)
}

func NewCommand(s *discordgo.Session, services service.All, logger logrus.FieldLogger) command {
	c := command{
		session:  s,
		services: services,
		log:      logger,
	}

	c.events = map[string]func(*discordgo.InteractionCreate, *payload) (entity.Payload, error){
		"":           c.InteractionCreate,
		eventConfirm: c.Confirm,
		eventCancel:  c.Cancel,
	}

	return c
}

func (c command) Name() string {
	return "ping"

}
func (c command) Description() string {
	return "ping"
}

func (c command) Options() []*discordgo.ApplicationCommandOption {
	return nil
}

func (c command) Cooldown() int {
	return 1
}

func (c command) Payload(key string) (*payload, error) {
	var payload payload
	return &payload, c.services.CustomID.Get(key, &payload)
}

func (c command) ExecEvent(i *discordgo.InteractionCreate, customID *discord.CustomID) (entity.Payload, error) {
	payload, err := c.Payload(customID.PayloadKey)
	if err != nil {
		return nil, err
	}

	execEvent, ok := c.events[customID.Event]
	if !ok {
		return nil, errors.ErrEventNotFound
	}

	payload.CustomID = customID
	return execEvent(i, payload)
}
