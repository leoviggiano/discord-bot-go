package dserr

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/constants/colors"
	"bot_test/embeds"
	"bot_test/service"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

var errMap = map[error]struct {
	ephemeral bool
	embed     func(*discordgo.Session, *discordgo.InteractionCreate) *discordgo.MessageEmbed
}{
	ErrUnauthorized: {ephemeral: true, embed: errUnauthorized},
}

func genericError(s *discordgo.Session) *discordgo.InteractionResponse {
	embed := embeds.DefaultBot(s)
	embed.Color = colors.Failure
	embed.Description = "Deu erro"

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
}

func Error(s *discordgo.Session, i *discordgo.InteractionCreate, services service.All, log logrus.FieldLogger, err error) *discordgo.InteractionResponse {
	mappedError, ok := errMap[err]
	if !ok {
		log.Error(err)
		return genericError(s)
	}

	var flags discordgo.MessageFlags
	if mappedError.ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{mappedError.embed(s, i)},
			Flags:  flags,
		},
	}
}
