package middleware

import (
	"github.com/bwmarrin/discordgo"

	"megumin/pkg/lib/discord"
	"megumin/pkg/listener/interaction/command"
	"megumin/pkg/service"
)

type Args struct {
	Session           discord.Session
	InteractionCreate *discordgo.InteractionCreate
	Command           command.Command
	Services          service.All
	Language          string
	Event             string
	CustomID          discord.CustomID
}

type Middleware interface {
	Exec(Args) (bool, error)
}

func All() []Middleware {
	return []Middleware{
		NewCanceller(),
		NewCooldown(),
		NewPermission(),
		NewUser(),
	}
}
