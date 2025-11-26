package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig `mapstructure:"server"`
	RedisConfig  `mapstructure:"redis"`
	MySQLConfig  `mapstructure:"mysql"`
	JwtConfig    `mapstructure:"jwt"`
}

var config Config

func init() {

	viper.SetConfigFile("./config.yaml")
	viper.ReadInConfig()
	viper.Unmarshal(&config)
	fmt.Println(config)
}
func Server() ServerConfig {
	return config.ServerConfig
}

func Redis() RedisConfig {
	return config.RedisConfig
}

func MySQL() MySQLConfig {
	return config.MySQLConfig
}

func Jwt() JwtConfig {
	return config.JwtConfig
}
