package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Emoji struct {
	ID       string
	Name     string
	Animated bool
}

func (e Emoji) String() string {
	if e.Name == "" {
		return ""
	}

	animated := ""
	if e.Animated {
		animated = "a"
	}

	return fmt.Sprintf("<%s:%s:%s>", animated, e.Name, e.ID)
}

func (e Emoji) Discord() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{
		ID:       e.ID,
		Name:     e.Name,
		Animated: e.Animated,
	}
}

func (e Emoji) ImageURL() string {
	if e.Animated {
		return fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.gif?size=256&quality=lossless", e.ID)
	}

	return fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.webp?size=256&quality=lossless", e.ID)
}

func NewEmoji(emojiStr string) Emoji {
	if emojiStr == "" {
		return Emoji{}
	}

	if !strings.HasPrefix(emojiStr, "<") || !strings.HasSuffix(emojiStr, ">") {
		return Emoji{
			Name: emojiStr,
		}
	}

	cleanedStr := strings.Trim(emojiStr, "<>")

	parts := strings.Split(cleanedStr, ":")
	if len(parts) != 3 || (parts[0] != "a" && parts[0] != "") {
		return Emoji{}
	}

	return Emoji{
		Animated: parts[0] == "a",
		Name:     parts[1],
		ID:       parts[2],
	}
}
