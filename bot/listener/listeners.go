package listener

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/service"
)

func All(session *discordgo.Session, services service.All, log logrus.FieldLogger) []interface{} {
	return []interface{}{
		NewReadyEvent(log).Handler,
		NewInteraction(session, services, log).Handler,
	}
}
