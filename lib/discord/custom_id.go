package discord

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	CustomIDSeparator = "<|>"

	CommandNameIndex = 0
	PayloadKeyIndex  = 1
	EventIndex       = 2

	CustomIDMaxLength = 100
)

type CustomID struct {
	PayloadKey  string
	CommandName string
	Event       string
}

type CustomIDOption func(*CustomID)

func NewCustomID(commandName string, options ...CustomIDOption) *CustomID {
	customID := &CustomID{
		CommandName: commandName,
	}

	for _, option := range options {
		option(customID)
	}

	return customID
}

func (c *CustomID) String() (string, error) {
	customID := make([]string, 3)
	customID[CommandNameIndex] = c.CommandName
	customID[PayloadKeyIndex] = c.PayloadKey
	customID[EventIndex] = c.Event

	fullCustomID := strings.Join(customID, CustomIDSeparator)
	if len(fullCustomID) > CustomIDMaxLength {
		return "", errors.New("custom id exceding max 100 chars")
	}

	return fullCustomID, nil
}

func CustomIDFromInteraction(i *discordgo.InteractionCreate) *CustomID {
	var customID, event string
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		customID = data.CustomID
		if len(data.Values) > 0 {
			event = data.Values[0]
		}

	case discordgo.InteractionModalSubmit:
		customID = i.ModalSubmitData().CustomID
	}

	if customID == "" {
		return nil
	}

	splitted := strings.Split(customID, CustomIDSeparator)

	if event == "" {
		event = splitted[EventIndex]
	}

	return &CustomID{
		CommandName: splitted[CommandNameIndex],
		PayloadKey:  splitted[PayloadKeyIndex],
		Event:       event,
	}
}

func WithPayloadKey(payloadKey string) CustomIDOption {
	return func(ci *CustomID) {
		ci.PayloadKey = payloadKey
	}
}
