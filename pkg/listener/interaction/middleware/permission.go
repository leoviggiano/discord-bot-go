package middleware

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/bot/embeds"
	"megumin/pkg/constants/colors"
)

type permission struct{}

func NewPermission() Middleware {
	return &permission{}
}

func (mw *permission) Exec(args Args) (next bool, err error) {
	if args.Command.Permission() == 0 {
		next = true
		return
	}

	// user := discord.UserFromInteraction(i)
	// if config.Admin(user.ID) {
	// 	next = true
	// 	return
	// }

	err = args.Session.InteractionRespond(args.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       colors.Red,
					Description: fmt.Sprintf("Você precisa ser um administrador para usar o comando `%s`", args.Command.Name()),
					Timestamp:   embeds.Timestamp(),
					Footer:      embeds.BotFooter(args.Session, ""),
				},
			},
		},
	})

	next = false
	return
}
