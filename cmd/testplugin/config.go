package main

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

type Config struct {
	Channel  string
	RPCHost  string
	RPCPort  int
	RPCToken string
}

func setDefault(conf *viper.Viper) {
	conf.SetDefault("RPC_HOST", "localhost")
	conf.SetDefault("RPC_PORT", 8001)
	conf.SetDefault("RPC_TOKEN", "")
	conf.SetDefault("CHANNEL", "")
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
	return &Config{
		Channel:  conf.GetString("CHANNEL"),
		RPCHost:  conf.GetString("RPC_HOST"),
		RPCPort:  conf.GetInt("RPC_PORT"),
		RPCToken: conf.GetString("RPC_TOKEN"),
	}, nil
}
