package dserr

import (
	"github.com/bwmarrin/discordgo"

	"bot_test/constants/colors"
	"bot_test/embeds"
)

func errCommandExpired(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	embed := embeds.DefaultBot(s)
	embed.Color = colors.Failure
	embed.Description = "Comando expirou"

	return embed
}

func errUnauthorized(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	embed := embeds.DefaultBot(s)
	embed.Color = colors.Failure
	embed.Description = "Você não está autorizado a fazer essa ação"

	return embed
}
