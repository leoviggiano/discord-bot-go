package middleware

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/bot/embeds"
	"megumin/pkg/constants/colors"
	"megumin/pkg/lib/discord"
	"megumin/pkg/listener/interaction/command"
)

// cooldown struct
/*
	cooldown {
		[command]: {
			[event]: {
				[userId]: timestamp
			}
		}
	}
*/

type cooldown struct {
	cooldowns map[string]map[string]map[string]int
}

func NewCooldown() Middleware {
	return &cooldown{
		cooldowns: make(map[string]map[string]map[string]int),
	}
}

func (mw *cooldown) addCommandToMap(cmd, event string) {
	mw.cooldowns[cmd] = make(map[string]map[string]int)
	mw.cooldowns[cmd][event] = make(map[string]int)
}

func (mw *cooldown) addUserToMap(userID, event string, cmd command.Command) {
	timestamp := time.Now().Add(time.Second * time.Duration(cmd.Cooldown()))

	if mw.cooldowns[cmd.Name()] == nil || mw.cooldowns[cmd.Name()][event] == nil {
		mw.addCommandToMap(cmd.Name(), event)
	}

	mw.cooldowns[cmd.Name()][event][userID] = int(timestamp.Unix())

	time.AfterFunc(time.Second*time.Duration(cmd.Cooldown()), func() {
		delete(mw.cooldowns[cmd.Name()][event], userID)
	})
}

func (mw *cooldown) Exec(args Args) (next bool, err error) {
	userID := discord.UserFromInteraction(args.InteractionCreate).ID
	timestamp := mw.cooldowns[args.Command.Name()][args.Event][userID]
	if timestamp == 0 {
		mw.addUserToMap(userID, args.Event, args.Command)
		next = true
		return
	}

	currentTime := int(time.Now().Unix())

	if currentTime > timestamp {
		mw.addUserToMap(userID, args.Event, args.Command)
		next = true
		return
	}

	timeLeft := timestamp - int(time.Now().Unix())

	err = args.Session.InteractionRespond(args.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       colors.Red,
					Title:       "Comando em Cooldown",
					Description: fmt.Sprintf("Por favor, aguarde mais %d segundos antes de usar este comando novamente", timeLeft),
					Timestamp:   embeds.Timestamp(),
					Footer:      embeds.BotFooter(args.Session, ""),
				},
			},
		},
	})

	return
}
