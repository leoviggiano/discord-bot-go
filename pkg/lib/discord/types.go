package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const Empty = "\u200b"

func EmptyField(inline ...bool) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   Empty,
		Value:  Empty,
		Inline: len(inline) > 0 && inline[0],
	}
}

var (
	ErrWithFollowUp = errors.New("follow up")
)

func ErrFollowUp(err error) error {
	return fmt.Errorf("%w: %w", ErrWithFollowUp, err)
}

func NewMessage() discordgo.InteractionResponseType {
	return discordgo.InteractionResponseChannelMessageWithSource
}

func UpdateMessage() discordgo.InteractionResponseType {
	return discordgo.InteractionResponseUpdateMessage
}
