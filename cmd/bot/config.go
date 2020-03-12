package main

import (
	"errors"
	"github.com/greboid/irc/database"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
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
	Debug    bool
	SASLAuth bool
	SASLUser string
	SASLPass string
	Plugins  []database.Plugin
}

func setDefault(conf *viper.Viper) {
	conf.SetDefault("SERVER", "")
	conf.SetDefault("PASSWORD", "")
	conf.SetDefault("TLS", true)
	conf.SetDefault("NICK", "")
	conf.SetDefault("WEB_PORT", 8000)
	conf.SetDefault("CHANNEL", "")
	conf.SetDefault("DB_PATH", "./data/db")
	conf.SetDefault("ADMIN_KEY", "ctwJTQ7HBdym3cns")
	conf.SetDefault("DEBUG", false)
	conf.SetDefault("SASL_AUTH", false)
	conf.SetDefault("SASL_USER", false)
	conf.SetDefault("SASL_PASS", false)
	conf.SetDefault("PLUGINS", "")
}

func getConfig(conf *viper.Viper) {
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(filepath.Join("$XDG_CONFIG", ".girc"))
	conf.AddConfigPath(filepath.Join("$XDG_CONFIG_HOME", ".girc"))
	conf.AddConfigPath(filepath.Join("$HOME", "config", ".girc"))
	conf.AddConfigPath(filepath.Join("$HOME", ".girc"))
	conf.AddConfigPath(".")
	if err := conf.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
		}
	}
}

func GetConfig() (*Config, error) {
	log.Print("Loading config")
	conf := viper.New()
	conf.AutomaticEnv()
	setDefault(conf)
	getConfig(conf)
	var plugins []database.Plugin
	for _, value := range strings.Split(conf.GetString("PLUGINS"), ",") {
		if len(value) == 0 {
			break
		}
		pluginString := strings.Split(value, "=")
		if len(pluginString) != 2 {
			return nil, errors.New("invalid plugin definition")
		}
		plugins = append(plugins, database.Plugin{Name: pluginString[0], Token: pluginString[1]})
	}
	return &Config{
		WebPort:  conf.GetInt("WEB_PORT"),
		Channel:  conf.GetString("CHANNEL"),
		Server:   conf.GetString("SERVER"),
		Password: conf.GetString("PASSWORD"),
		Nickname: conf.GetString("NICK"),
		DBPath:   conf.GetString("DB_PATH"),
		AdminKey: conf.GetString("ADMIN_KEY"),
		TLS:      conf.GetBool("TLS"),
		Debug:    conf.GetBool("DEBUG"),
		SASLAuth: conf.GetBool("SASL_AUTH"),
		SASLUser: conf.GetString("SASL_USER"),
		SASLPass: conf.GetString("SASL_PASS"),
		Plugins:  plugins,
	}, nil
}
