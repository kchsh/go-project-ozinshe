package config

import "time"

var Config *MapConfig

type MapConfig struct {
	AppPort            string        `mapstructure:"APP_PORT"`
	DbConnectionString string        `mapstructure:"DB_CONNECTION_STRING"`
	JwtSecretKey       string        `mapstructure:"JWT_SECRET_KEY"`
	JwtExpiresIn       time.Duration `mapstructure:"JWT_EXPIRE_DURATION"`
}
