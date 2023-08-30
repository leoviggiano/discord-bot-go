package embeds

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"bot_test/constants/colors"
)

const Empty = "\u200b"

func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func Default() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: colors.Default,
	}
}

func DefaultBot(s *discordgo.Session) *discordgo.MessageEmbed {
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

func BotFooter(s *discordgo.Session, text string) *discordgo.MessageEmbedFooter {
	if text == "" {
		text = s.State.User.Username
	}

	return &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: s.State.User.AvatarURL(""),
	}
}

func BotAuthor(s *discordgo.Session) *discordgo.MessageEmbedAuthor {
	return &discordgo.MessageEmbedAuthor{
		Name:    s.State.User.Username,
		IconURL: s.State.User.AvatarURL(""),
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
