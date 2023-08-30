package middleware

import (
	"github.com/bwmarrin/discordgo"

	"bot_test/bot/command"
	"bot_test/service"
)

type Middleware interface {
	Exec(*discordgo.Session, *discordgo.InteractionCreate, command.Command, service.All) (next bool, err error)
}

func All() []Middleware {
	return []Middleware{
		NewCooldown(),
	}
}
