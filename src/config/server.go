package config

import "strconv"

type ServerConfig struct {
	Port            int    `mapstructure:"port"`
	ShortCodeLength int    `mapstructure:"short_code_length"`
	Domain          string `mapstructure:"domain"`
}

func (s ServerConfig) PortStr() string {
	return ":" + strconv.Itoa(s.Port)
}
