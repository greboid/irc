package main

import "github.com/spf13/viper"

type Config struct {
	WebPort  int
	channel  string
	server   string
	password string
	nickname string
}

func getConfig() Config {
	viper.SetDefault("SERVER", "")
	viper.SetDefault("PASSWORD", "")
	viper.SetDefault("NICK", "")
	viper.SetDefault("WEB_PORT", 8000)
	viper.SetDefault("CHANNEL", "")
	viper.AutomaticEnv()
	return Config{
		WebPort:  viper.GetInt("WEB_PORT"),
		channel:  viper.GetString("CHANNEL"),
		server:   viper.GetString("SERVER"),
		password: viper.GetString("PASSWORD"),
		nickname: viper.GetString("NICK"),
	}
}
