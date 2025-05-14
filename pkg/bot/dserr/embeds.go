package dserr

import (
	"github.com/bwmarrin/discordgo"

	"megumin/pkg/bot/embeds"
	"megumin/pkg/constants/colors"
)

func (d dsErr) errCommandExpired(i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	embed := embeds.DefaultBot(d.session)
	embed.Color = colors.Failure
	embed.Description = "Comando expirado"

	return embed
}
