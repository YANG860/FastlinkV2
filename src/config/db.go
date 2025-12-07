package config

import "time"

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr                   string        `mapstructure:"addr"`
	Password               string        `mapstructure:"password"`
	DB                     int           `mapstructure:"db"`
	RefreshTokenTTL        time.Duration `mapstructure:"refreshTokenTTL"`
	LinkTTL                time.Duration `mapstructure:"linkTTL"`
	ClickWriteBackInterval time.Duration `mapstructure:"clickWriteBackInterval"`
}
