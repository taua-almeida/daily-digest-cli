package internal

import "github.com/BurntSushi/toml"

type PrintStyle struct {
	Style string `toml:"style"`
}

type RateConfig struct {
	RateType  string `toml:"type"`
	RateLimit int    `toml:"limit"`
}

type Config struct {
	Title      string
	Print      PrintStyle `toml:"print"`
	RateConfig RateConfig `toml:"rate"`
}

func ReadConfig(configPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
