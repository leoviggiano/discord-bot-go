package listener

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type ReadyEvent struct {
	log logrus.FieldLogger
}

func NewReadyEvent(log logrus.FieldLogger) *ReadyEvent {
	return &ReadyEvent{log}
}

func (l *ReadyEvent) Handler(s *discordgo.Session, e *discordgo.Ready) {
	l.log.Info("[ListenerReady] Bot session is ready")
	l.log.Infof("[ListenerReady] Logged in as %s", e.User.String())

	presence := presence{
		Game:   "TESTE",
		Status: "dnd",
	}

	s.UpdateStatusComplex(presence.updatePresence())
}

type presence struct {
	Game   string
	Status string
}

func (p presence) updatePresence() discordgo.UpdateStatusData {
	return discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: p.Game,
				Type: discordgo.ActivityTypeGame,
			},
		},
		Status: p.Status,
	}
}
