package middleware

import (
	"time"

	"megumin/pkg/constants"
	"megumin/pkg/lib/discord"
	"megumin/pkg/lib/snowflake"
)

// canceller struct
/*
	canceller {
		[command]: {
			[userId]: snowflake(custom id)
		}
	}
*/

type canceller struct {
	commands map[string]map[string]snowflake.Snowflake
}

func NewCanceller() Middleware {
	return &canceller{
		commands: make(map[string]map[string]snowflake.Snowflake),
	}
}

func (mw *canceller) addSnowflakeToUser(command string, userID string, currentSnowflake snowflake.Snowflake) {
	if _, ok := mw.commands[command]; !ok {
		mw.commands[command] = make(map[string]snowflake.Snowflake)
	}

	mw.commands[command][userID] = currentSnowflake

	time.AfterFunc(constants.CommandExpireTime, func() {
		if snowflake, ok := mw.commands[command][userID]; ok && snowflake == currentSnowflake {
			delete(mw.commands[command], userID)
		}
	})
}

func (mw *canceller) Exec(args Args) (next bool, err error) {
	cmd := args.Command
	user := discord.UserFromInteraction(args.InteractionCreate)

	currentSnowflake, err := snowflake.FromString(args.CustomID.Key())
	if err != nil {
		return false, err
	}

	command, ok := mw.commands[cmd.Name()]
	if !ok {
		mw.addSnowflakeToUser(cmd.Name(), user.ID, currentSnowflake)
		return true, nil
	}

	lastSnowflake, ok := command[user.ID]
	if !ok {
		mw.addSnowflakeToUser(cmd.Name(), user.ID, currentSnowflake)
		return true, nil
	}

	if lastSnowflake > currentSnowflake {
		return false, args.Session.InteractionRespond(args.InteractionCreate.Interaction, discord.InteractionNoResponse())
	}

	mw.addSnowflakeToUser(cmd.Name(), user.ID, currentSnowflake)
	return true, nil
}
