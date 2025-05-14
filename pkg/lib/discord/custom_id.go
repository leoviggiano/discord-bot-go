package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	CustomIDDelimiter = "<|>"
	EventDelimiter    = "e:"
)

type CustomID string

func (c CustomID) String() string {
	return string(c)
}

func (c CustomID) WithValue(values ...any) CustomID {
	customID := c.Key()
	for _, value := range values {
		customID += fmt.Sprintf("%s%v", CustomIDDelimiter, value)
	}

	return CustomID(customID)
}

func (c CustomID) WithEvent(values ...any) CustomID {
	customID := c.Key()
	for _, value := range values {
		customID += fmt.Sprintf("%s%s%v", CustomIDDelimiter, EventDelimiter, value)
	}

	return CustomID(customID)
}

func (c CustomID) Event() string {
	splitted := strings.Split(string(c), EventDelimiter)
	if len(splitted) > 1 {
		return splitted[1]
	}

	return ""
}

func (c CustomID) Values() []string {
	return strings.Split(string(c), CustomIDDelimiter)[1:]
}

func (c CustomID) Key() string {
	return strings.Split(string(c), CustomIDDelimiter)[0]
}

func PayloadKeyFromInteraction(i *discordgo.InteractionCreate) string {
	var customID string
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		customID = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		customID = i.MessageComponentData().CustomID
	case discordgo.InteractionModalSubmit:
		customID = i.ModalSubmitData().CustomID
	}

	return customID
}
