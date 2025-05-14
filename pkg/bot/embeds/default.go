package embeds

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/constants/colors"
	"megumin/pkg/lib/discord"
)

func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func Default() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: colors.Default,
	}
}

func DefaultBot(s discord.Session) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     colors.Default,
		Timestamp: Timestamp(),
		Footer:    BotFooter(s, ""),
		Author:    BotAuthor(s),
	}
}

func BotProfileAvatar(s *discordgo.Session) string {
	return s.State.User.AvatarURL("")
}

func BotFooter(s discord.Session, text string) *discordgo.MessageEmbedFooter {
	if text == "" {
		text = s.Username()
	}

	return &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: s.AvatarURL(),
	}
}

func BotAuthor(s discord.Session) *discordgo.MessageEmbedAuthor {
	return &discordgo.MessageEmbedAuthor{
		Name:    s.Username(),
		IconURL: s.AvatarURL(),
	}
}

func UserAuthor(u *discordgo.User) *discordgo.MessageEmbedAuthor {
	return &discordgo.MessageEmbedAuthor{
		Name:    u.Username,
		IconURL: u.AvatarURL(""),
	}
}

func EmptyField(inline ...bool) *discordgo.MessageEmbedField {
	isInline := false
	if len(inline) > 0 {
		isInline = inline[0]
	}

	return &discordgo.MessageEmbedField{
		Name:   "\u200b",
		Value:  "\u200b",
		Inline: isInline,
	}
}

func EmptyComponent() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{}
}
