package ready

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type ReadyEvent struct {
	log *slog.Logger
}

func ReadyEventListener(log *slog.Logger) *ReadyEvent {
	return &ReadyEvent{log}
}

func (l *ReadyEvent) Handler(s *discordgo.Session, e *discordgo.Ready) {
	l.log.Info("[ListenerReady] Bot session is ready")
	l.log.Info("[ListenerReady] Logged in as", "user", e.User.String())

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
