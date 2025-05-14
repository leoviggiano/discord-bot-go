package context

import (
	"slices"

	"megumin/pkg/lib/discord"
)

type Payload struct {
	CustomID   discord.CustomID // Payload Key
	Command    string           // Command that will be executed
	Event      string           // Event that will be executed
	ForceEvent string           // Force event to be executed
	Reacters   []string         // Who can react
	Options    any              // Especific command struct
}

func (p *Payload) React(userID string, ephemeral bool) bool {
	if len(p.Reacters) == 0 || ephemeral {
		return true
	}

	return slices.Contains(p.Reacters, userID)
}
