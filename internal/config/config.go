package config

import (
	"fmt"
	"mistar-be-go/pkg/appinfo"
	"mistar-be-go/pkg/logger"
	"mistar-be-go/pkg/pgsql"
	"mistar-be-go/pkg/token"
	"os"

	"github.com/pelletier/go-toml"
)

type API struct {
	Host     string `toml:"host"`
	RESTPort int    `toml:"rest_port"`
}

type OAuth struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_uri"`
}

type DevSettings struct {
	Auth bool `toml:"auth"`
}

type Config struct {
	API         API             `toml:"api"`
	AppInfo     appinfo.Info    `toml:"app_info"`
	Logger      logger.Config   `toml:"logger"`
	PostgreSQL  pgsql.Config    `toml:"postgres"`
	OAuth       OAuth           `toml:"oauth"`
	JWT         token.Config    `toml:"jwt"`
	DevSettings DevSettings     `toml:"dev"`
}

func LoadEnvFromFile(path string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("error open config file: %w", err)
	}
	err = toml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("error parsing toml: %w", err)
	}

	return cfg, nil
}