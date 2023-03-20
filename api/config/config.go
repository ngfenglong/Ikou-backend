package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port            int    `mapstructure:"PORT"`
	FrontendAddress string `mapstructure:"FRONTEND_ADDRESS"`
	DbDriver        string `mapstructure:"DB_DRIVER"`
	DbSource        string `mapstructure:"DB_SOURCE"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
