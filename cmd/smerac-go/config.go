package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Discord struct {
	Token   string `mapstructure:"TOKEN"`
	GuildId string `mapstructure:"GUILDID"`
}

type Google struct {
	Token string `mapstructure:"TOKEN"`
}

type Role struct {
	Id   string `mapstructure:"ID"`
	Name string `mapstructure:"NAME"`
}

type NamedDays struct {
	Monday    string `mapstructure:"MONDAY"`
	Tuesday   string `mapstructure:"TUESDAY"`
	Wednesday string `mapstructure:"WEDNESDAY"`
	Thursday  string `mapstructure:"THURSDAY"`
	Friday    string `mapstructure:"FRIDAY"`
	Saturday  string `mapstructure:"SATURDAY"`
	Sunday    string `mapstructure:"SUNDAY"`
}

type Calendar struct {
	Id                string    `mapstructure:"ID"`
	ChannelId         string    `mapstructure:"CHANNELID"`
	Name              string    `mapstructure:"NAME"`
	TimeBetweenChecks int       `mapstructure:"TIME"`
	NamedDays         NamedDays `mapstructure:"NAMEDDAYS"`
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

	for index, calendar := range config.Calendars {
		if calendar.NamedDays.Monday == "" {
			config.Calendars[index].NamedDays.Monday = "Monday"
		}
		if calendar.NamedDays.Tuesday == "" {
			config.Calendars[index].NamedDays.Tuesday = "Tuesday"
		}
		if calendar.NamedDays.Wednesday == "" {
			config.Calendars[index].NamedDays.Wednesday = "Wednesday"
		}
		if calendar.NamedDays.Thursday == "" {
			config.Calendars[index].NamedDays.Thursday = "Thursday"
		}
		if calendar.NamedDays.Friday == "" {
			config.Calendars[index].NamedDays.Friday = "Friday"
		}
		if calendar.NamedDays.Saturday == "" {
			config.Calendars[index].NamedDays.Saturday = "Saturday"
		}
		if calendar.NamedDays.Sunday == "" {
			config.Calendars[index].NamedDays.Sunday = "Sunday"
		}
	}

	return config, err
}
