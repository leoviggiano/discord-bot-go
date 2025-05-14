package ping

import (
	"github.com/bwmarrin/discordgo"

	"megumin/pkg/constants/events"
	"megumin/pkg/listener/interaction/context"
)

func (c command) InteractionCreate(ctx *context.Context) (string, error) {
	ctx.Session.InteractionRespond(ctx.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})

	return events.CommandEnded, nil
}
