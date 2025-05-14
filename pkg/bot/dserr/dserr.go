package dserr

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/bot/embeds"
	"megumin/pkg/constants/colors"
	"megumin/pkg/constants/errors"
	"megumin/pkg/lib/discord"
)

type ErrorHandler interface {
	Error(i *discordgo.InteractionCreate, err error) *discordgo.InteractionResponse
}

type dsErr struct {
	session       discord.Session
	log           *slog.Logger
	errMap        map[error]func(*discordgo.InteractionCreate) *discordgo.MessageEmbed
	epheremalErrs map[error]bool
}

func NewErrorHandler(s discord.Session, log *slog.Logger) ErrorHandler {
	dsErr := &dsErr{
		session: s,
		log:     log,
	}

	errMap := map[error]func(*discordgo.InteractionCreate) *discordgo.MessageEmbed{
		errors.ErrCommandExpired: dsErr.errCommandExpired,
	}

	ephemeralErrs := map[error]bool{}

	dsErr.errMap = errMap
	dsErr.epheremalErrs = ephemeralErrs

	return dsErr
}

func (d dsErr) genericError() *discordgo.InteractionResponse {
	embed := embeds.DefaultBot(d.session)
	embed.Color = colors.Failure
	embed.Description = "Erro desconhecido"

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
}

func (d dsErr) Error(i *discordgo.InteractionCreate, err error) *discordgo.InteractionResponse {
	var flags discordgo.MessageFlags
	if d.epheremalErrs[err] {
		flags = discordgo.MessageFlagsEphemeral
	}

	if errEmbed, ok := d.errMap[err]; ok {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{errEmbed(i)},
				Flags:  flags,
			},
		}
	}

	d.log.Error("Unknown error", "error", err)

	return d.genericError()
}
