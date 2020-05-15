package services

import (
	"github.com/BurntSushi/toml"
)

type ConfigService struct {
	AppName string             `toml:"app_name"`
	Version string             `toml:"version"`
	FeedURL string             `toml:"feed_url"`
	MySQL   mysqlConfig        `toml:"mysql"`
	HTTP    httpEndpointConfig `toml:"http"`
}

type mysqlConfig struct {
	Address  string
	Database string
	User     string
	Password string
	PoolSize int `toml:"pool_size"`
}

type httpEndpointConfig struct {
	Address   string `toml:"address"`
	RateLimit int    `toml:"rate_limit"`
}

func (c *ConfigService) Load(file string) error {
	if _, err := toml.DecodeFile(file, &c); err != nil {
		return err
	}

	return nil
}
