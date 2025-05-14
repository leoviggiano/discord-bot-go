package discord

import (
	"github.com/bwmarrin/discordgo"
)

func UserFromInteraction(i *discordgo.InteractionCreate) *discordgo.User {
	if i.User != nil {
		return i.User
	}

	return i.Member.User
}

func MaxOptions(max, currentLength int) int {
	if currentLength < max {
		return currentLength
	}
	return max
}

func InteractionNoResponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}
}

func ValuesFromModal(data discordgo.ModalSubmitInteractionData) map[string]string {
	values := make(map[string]string)

	for _, d := range data.Components {
		input := d.(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
		values[input.CustomID] = input.Value
	}

	return values
}
