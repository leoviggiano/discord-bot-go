package middleware

import (
	"megumin/pkg/lib/discord"
)

type user struct{}

func NewUser() Middleware {
	return &user{}
}

func (mw *user) Exec(args Args) (next bool, err error) {
	user := discord.UserFromInteraction(args.InteractionCreate)
	exists := args.Services.User.Exists(user.ID)

	if !exists {
		err = args.Services.User.Create(user.ID)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
