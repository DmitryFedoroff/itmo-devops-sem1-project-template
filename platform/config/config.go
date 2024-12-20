package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const DefaultDateFormat = "2006-01-02"

type ServerSettings struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseSettings struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type Settings struct {
	Server   ServerSettings   `mapstructure:"server"`
	Database DatabaseSettings `mapstructure:"database"`
}

func Load(configPath string) (Settings, error) {
	log.Printf("Loading configuration from '%s'\n", configPath)

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Settings{}, err
	}

	var cfg Settings
	if err := viper.Unmarshal(&cfg); err != nil {
		return Settings{}, err
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}
