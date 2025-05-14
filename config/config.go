package config

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	Production  = "production"
	Development = "development"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}

	viper.AutomaticEnv()
}

func Environment() string {
	switch viper.GetString("environment") {
	case "development":
		return Development
	default:
		return Production
	}
}

func DatabaseConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("db_user"),
		viper.GetString("db_password"),
		viper.GetString("db_host"),
		viper.GetString("db_port"),
		viper.GetString("db_name"),
	)
}

func DiscordToken() string {
	return viper.GetString("discord_token")
}

func GuildID() string {
	if Environment() == Production {
		return ""
	}

	return viper.GetString("discord_guild_id")
}

func AdminGuildID() string {
	return viper.GetString("discord_admin_guild_id")
}

func AppID() string {
	return viper.GetString("discord_app_id")
}
