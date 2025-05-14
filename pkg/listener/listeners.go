package listener

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/lib/sql"
	"megumin/pkg/listener/interaction"
	"megumin/pkg/listener/ready"
)

func All(session *discordgo.Session, db sql.Connection, log *slog.Logger) ([]interface{}, error) {
	interactionListener, err := interaction.NewListener(session, db, log)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		interactionListener.Handler,
		ready.ReadyEventListener(log).Handler,
	}, nil
}
