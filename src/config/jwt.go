package config

import "time"

type JwtConfig struct {
	JwtKey               string        `mapstructure:"jwtKey"`
	AccessTokenTTL       time.Duration `mapstructure:"accessTokenTTL"`
	RefreshTokenTTL      time.Duration `mapstructure:"refreshTokenTTL"`
	RefreshTokenIDLength int           `mapstructure:"refreshTokenIDLength"`
}
