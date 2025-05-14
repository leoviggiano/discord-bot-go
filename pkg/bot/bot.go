package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"megumin/pkg/lib/sql"
	"megumin/pkg/listener"
)

type Discord struct {
	Session *discordgo.Session
	Logger  *slog.Logger
}

const (
	Intents = discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsDirectMessageReactions
)

func StartBot(token string, db sql.Connection, log *slog.Logger) error {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	session.Identify.Intents = Intents

	discord := &Discord{
		Session: session,
		Logger:  log,
	}

	discord.registerListeners(db)

	if err := session.Open(); err != nil {
		return err
	}

	log.Info("[Discord] Bot connected with success")
	return nil
}

func (d *Discord) registerListeners(db sql.Connection) {
	listeners, err := listener.All(d.Session, db, d.Logger)
	if err != nil {
		d.Logger.Error("Error registering listeners", "error", err)
		return
	}

	for _, listener := range listeners {
		d.Session.AddHandler(listener)
	}

	d.Logger.Info("[RegisterListeners] All listeners registered")
}
