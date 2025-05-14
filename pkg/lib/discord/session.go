//go:generate mockgen -source=session.go -destination=../../mocks/session.go -package=mocks
package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Session interface {
	ApplicationCommands(appID, guildID string) (cmd []*discordgo.ApplicationCommand, err error)
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand) (ccmd *discordgo.ApplicationCommand, err error)
	ApplicationCommandBulkOverwrite(appID string, guildID string, commands []*discordgo.ApplicationCommand) (createdCommands []*discordgo.ApplicationCommand, err error)
	ApplicationCommandDelete(appID, guildID, cmdID string, options ...discordgo.RequestOption) error

	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error
	InteractionResponseEdit(interaction *discordgo.Interaction, newresp *discordgo.WebhookEdit) (*discordgo.Message, error)
	InteractionResponseDelete(interaction *discordgo.Interaction) error

	FollowupMessageCreate(interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error)
	FollowupMessageEdit(interaction *discordgo.Interaction, messageID string, data *discordgo.WebhookEdit) (*discordgo.Message, error)
	FollowupMessageDelete(interaction *discordgo.Interaction, messageID string) error

	AvatarURL() string
	Username() string
}

type discord struct {
	session *discordgo.Session
}

func NewSession(s *discordgo.Session) Session {
	return &discord{session: s}
}

func (d *discord) ApplicationCommands(appID, guildID string) (cmd []*discordgo.ApplicationCommand, err error) {
	return d.session.ApplicationCommands(appID, guildID)
}

func (d *discord) ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand) (ccmd *discordgo.ApplicationCommand, err error) {
	return d.session.ApplicationCommandCreate(appID, guildID, cmd)
}

func (d *discord) ApplicationCommandBulkOverwrite(appID string, guildID string, commands []*discordgo.ApplicationCommand) (createdCommands []*discordgo.ApplicationCommand, err error) {
	return d.session.ApplicationCommandBulkOverwrite(appID, guildID, commands)
}

func (d *discord) ApplicationCommandDelete(appID, guildID, cmdID string, options ...discordgo.RequestOption) error {
	return d.session.ApplicationCommandDelete(appID, guildID, cmdID, options...)
}

func (d *discord) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	if resp.Type == 0 {
		resp.Type = discordgo.InteractionResponseChannelMessageWithSource
	}

	if resp.Type == discordgo.InteractionResponseUpdateMessage &&
		interaction.Type == discordgo.InteractionApplicationCommand {
		var (
			content    *string
			components *[]discordgo.MessageComponent
			embeds     *[]*discordgo.MessageEmbed
		)

		if data := resp.Data; data != nil {
			if data.Content != "" {
				content = &data.Content
			}

			if len(data.Components) > 0 {
				components = &data.Components
			}

			if len(data.Embeds) > 0 {
				embeds = &data.Embeds
			}
		}

		_, err := d.session.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
			Content:         content,
			Components:      components,
			Embeds:          embeds,
			Files:           resp.Data.Files,
			Attachments:     resp.Data.Attachments,
			AllowedMentions: resp.Data.AllowedMentions,
		})
		return err
	}

	return d.session.InteractionRespond(interaction, resp)
}

func (d *discord) InteractionResponseEdit(interaction *discordgo.Interaction, newresp *discordgo.WebhookEdit) (*discordgo.Message, error) {
	return d.session.InteractionResponseEdit(interaction, newresp)
}

func (d *discord) InteractionResponseDelete(interaction *discordgo.Interaction) error {
	return d.session.InteractionResponseDelete(interaction)
}

func (d *discord) FollowupMessageCreate(interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	return d.session.FollowupMessageCreate(interaction, wait, data)
}

func (d *discord) FollowupMessageEdit(interaction *discordgo.Interaction, messageID string, data *discordgo.WebhookEdit) (*discordgo.Message, error) {
	return d.session.FollowupMessageEdit(interaction, messageID, data)
}

func (d *discord) FollowupMessageDelete(interaction *discordgo.Interaction, messageID string) error {
	return d.session.FollowupMessageDelete(interaction, messageID)
}

func (d *discord) AvatarURL() string {
	return d.session.State.User.AvatarURL("")
}

func (d *discord) Username() string {
	return d.session.State.User.Username
}
