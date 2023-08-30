package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"bot_test/bot/listener"
	"bot_test/config"
	"bot_test/service"
)

type Discord struct {
	Session  *discordgo.Session
	Services service.All
	Log      logrus.FieldLogger
}

const (
	Intents = discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsDirectMessageReactions
)

func StartBot(services service.All, log logrus.FieldLogger) error {
	session, err := discordgo.New("Bot " + config.GetToken())
	if err != nil {
		return err
	}

	session.Identify.Intents = Intents

	discord := &Discord{
		Session:  session,
		Services: services,
		Log:      log,
	}

	discord.registerListeners()

	if err := session.Open(); err != nil {
		return err
	}

	log.Info("[Discord] Bot connected with success")
	return nil
}

func (d *Discord) registerListeners() {
	for _, listener := range listener.All(d.Session, d.Services, d.Log) {
		d.Session.AddHandler(listener)
	}

	d.Log.Info("[RegisterListeners] All listeners registered")
}
