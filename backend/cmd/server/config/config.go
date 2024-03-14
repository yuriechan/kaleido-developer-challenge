package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config stores configuration extracted from environmental variables by using:
// https://github.com/kelseyhightower/envconfig
type Config struct {
	FireflyBaseUrl string `envconfig:"FIREFLY_BASE_URL" required:"true"`
	MysqlPassword  string `envconfig:"MYSQL_PASSWORD" required:"true"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("envconfig.Process: %w", err)
	}
	return &c, nil
}
