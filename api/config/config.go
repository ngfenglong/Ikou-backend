package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port            int    `mapstructure:"PORT"`
	FrontendAddress string `mapstructure:"FRONTEND_ADDRESS"`
	DbDriver        string `mapstructure:"DB_DRIVER"`
	DbSource        string `mapstructure:"DB_SOURCE"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
}

var C *Config

func LoadConfig() error {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
