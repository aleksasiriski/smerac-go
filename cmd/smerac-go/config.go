package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Discord struct {
	Token string `mapstructure:"TOKEN"`
}

type Google struct {
	Token string `mapstructure:"TOKEN"`
}

type Role struct {
	Id   string `mapstructure:"ID"`
	Name string `mapstructure:"NAME"`
}

type Calendar struct {
	Id                string `mapstructure:"ID"`
	ChannelId         string `mapstructure:"CHANNELID"`
	Name              string `mapstructure:"NAME"`
	TimeBetweenChecks int    `mapstructure:"TIME"`
}

type Config struct {
	Discord   Discord    `mapstructure:"DISCORD"`
	Google    Google     `mapstructure:"GOOGLE"`
	Roles     []Role     `mapstructure:"ROLES"`
	Calendars []Calendar `mapstructure:"CALENDARS"`
}

func LoadConfig(path string) (Config, error) {
	config := Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName("smerac")
	viper.SetConfigType("yaml")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("failed parsing config: %w", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed unmarshaling config: %w", err)
	}

	if config.Discord.Token == "" {
		return config, fmt.Errorf("discord token can't be empty")
	}

	return config, err
}
