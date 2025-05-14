package context

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/lib/discord"
)

type Context struct {
	Payload           *Payload
	CustomID          discord.CustomID
	Session           discord.Session
	InteractionCreate *discordgo.InteractionCreate
	Log               *slog.Logger
	User              *discordgo.User
}
