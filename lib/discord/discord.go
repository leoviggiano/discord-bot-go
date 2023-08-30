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

func ButtonsActionRow(customID *CustomID, buttons []Button) (discordgo.ActionsRow, error) {
	discordButtons := make([]discordgo.MessageComponent, len(buttons))
	for i, button := range buttons {
		newCustomID := &CustomID{
			PayloadKey:  customID.PayloadKey,
			CommandName: customID.CommandName,
			Event:       button.Value,
		}

		customIDString, err := newCustomID.String()
		if err != nil {
			return discordgo.ActionsRow{}, err
		}

		discordButtons[i] = discordgo.Button{
			Label:    button.Label,
			Style:    button.Style,
			Disabled: button.Disabled,
			Emoji:    button.Emoji,
			CustomID: customIDString,
		}
	}

	return discordgo.ActionsRow{
		Components: discordButtons,
	}, nil
}
