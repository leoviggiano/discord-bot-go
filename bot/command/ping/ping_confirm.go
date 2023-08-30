package ping

import (
	"github.com/bwmarrin/discordgo"

	"bot_test/embeds"
	"bot_test/entity"
	"bot_test/lib/discord"
)

func (c command) PingConfirm(i *discordgo.InteractionCreate, payload *payload) (entity.Payload, error) {
	embed := embeds.DefaultBot(c.session)
	embed.Description = "ping confirm"

	buttons := []discord.Button{
		{Value: eventConfirm, Label: "Confirmar", Style: discordgo.PrimaryButton},
		{Value: eventCancel, Label: "Cancelar", Style: discordgo.DangerButton},
	}

	buttonComponent, err := discord.ButtonsActionRow(payload.CustomID, buttons)
	if err != nil {
		return payload, err
	}

	err = c.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discord.NewMessage(),
		Data: &discordgo.InteractionResponseData{
			Embeds: discord.Embeds(embed),
			Components: []discordgo.MessageComponent{
				buttonComponent,
			},
		},
	})

	return payload, err
}
