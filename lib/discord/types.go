package discord

import "github.com/bwmarrin/discordgo"

func NewMessage() discordgo.InteractionResponseType {
	return discordgo.InteractionResponseChannelMessageWithSource
}

func UpdateMessage() discordgo.InteractionResponseType {
	return discordgo.InteractionResponseUpdateMessage
}

func Embeds(embeds ...*discordgo.MessageEmbed) []*discordgo.MessageEmbed {
	messageEmbeds := make([]*discordgo.MessageEmbed, 0, len(embeds))
	messageEmbeds = append(messageEmbeds, embeds...)

	return messageEmbeds
}

type Button struct {
	Value    string
	Label    string
	Style    discordgo.ButtonStyle
	Disabled bool
	Emoji    discordgo.ComponentEmoji
}
