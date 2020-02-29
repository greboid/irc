package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	WebPort  int
	Channel  string
	Server   string
	Password string
	Nickname string
	DBPath   string
	AdminKey string
	TLS      bool
}

func GetConfig() *Config {
	log.Print("Loading config")
	viper.SetDefault("PASSWORD", "")
	viper.SetDefault("NICK", "")
	viper.SetDefault("TLS", true)
	viper.SetDefault("WEB_PORT", 8000)
	viper.SetDefault("CHANNEL", "")
	viper.SetDefault("DB_PATH", "./data/db")
	viper.SetDefault("ADMIN_KEY", "ctwJTQ7HBdym3cns")
	viper.AutomaticEnv()
	log.Print("Returning config")
	return &Config{
		WebPort:  viper.GetInt("WEB_PORT"),
		Channel:  viper.GetString("CHANNEL"),
		Server:   viper.GetString("SERVER"),
		Password: viper.GetString("PASSWORD"),
		Nickname: viper.GetString("NICK"),
		DBPath:   viper.GetString("DB_PATH"),
		AdminKey: viper.GetString("ADMIN_KEY"),
		TLS:      viper.GetBool("TLS"),
	}
}
