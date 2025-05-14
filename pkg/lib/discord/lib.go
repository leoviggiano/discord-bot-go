package discord

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

type SelectMenu struct {
	CustomID    string                       `json:"custom_id,omitempty"`
	Placeholder string                       `json:"placeholder"`
	MinValues   int                          `json:"min_values"`
	MaxValues   int                          `json:"max_values,omitempty"`
	Options     []discordgo.SelectMenuOption `json:"options"`
}

func (SelectMenu) Type() discordgo.ComponentType {
	return discordgo.SelectMenuComponent
}

func (m SelectMenu) MarshalJSON() ([]byte, error) {
	type selectMenu SelectMenu

	return json.Marshal(struct {
		selectMenu
		Type discordgo.ComponentType `json:"type"`
	}{
		selectMenu: selectMenu(m),
		Type:       m.Type(),
	})
}
