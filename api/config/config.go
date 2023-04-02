package config

import (
	"fmt"

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

func LoadConfig(path string) error {
	viper.AddConfigPath(path)

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	fmt.Println("Port:", C.Port)
	fmt.Println("Frontend Address:", C.FrontendAddress)
	fmt.Println("DB Driver:", C.DbDriver)
	fmt.Println("DB Source:", C.DbSource)
	fmt.Println("Server Address:", C.ServerAddress)

	return nil
}
