package middleware

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"bot_test/bot/command"
	"bot_test/constants/colors"
	"bot_test/embeds"
	"bot_test/lib/discord"
	"bot_test/service"
)

// cooldown struct
/*
	cooldown {
		[command]: {
			[userId]: timestamp
		}
	}
*/

type cooldown struct {
	cooldowns map[string]map[string]int
}

func NewCooldown() Middleware {
	return &cooldown{
		cooldowns: make(map[string]map[string]int),
	}
}

func (mw *cooldown) addCommandToMap(cmd string) {
	mw.cooldowns[cmd] = make(map[string]int)
}

func (mw *cooldown) addUserToMap(userID string, cmd command.Command) {
	timestamp := time.Now().Add(time.Second * time.Duration(cmd.Cooldown()))

	if mw.cooldowns[cmd.Name()] == nil {
		mw.addCommandToMap(cmd.Name())
	}

	mw.cooldowns[cmd.Name()][userID] = int(timestamp.Unix())
}

func (mw *cooldown) Exec(s *discordgo.Session, i *discordgo.InteractionCreate, cmd command.Command, services service.All) (next bool, err error) {
	userID := discord.UserFromInteraction(i).ID
	timestamp := mw.cooldowns[cmd.Name()][userID]
	if timestamp == 0 {
		mw.addUserToMap(userID, cmd)
		next = true
		return
	}

	currentTime := int(time.Now().Unix())

	if currentTime > timestamp {
		mw.addUserToMap(userID, cmd)
		next = true
		return
	}

	timeLeft := timestamp - int(time.Now().Unix())

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       colors.Red,
					Title:       "Comando em cooldown",
					Description: fmt.Sprintf("Espere mais %d segundos antes de usar o comando %s", timeLeft, cmd.Name()),
					Timestamp:   embeds.Timestamp(),
					Footer:      embeds.BotFooter(s, ""),
				},
			},
		},
	})

	return
}
