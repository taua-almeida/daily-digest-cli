package internal

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

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
	validateConfigs(&config)
	return &config, nil
}

func validateConfigs(config *Config) error {
	if config.RateConfig.RateType != "" {
		err := isValidRateType(config.RateConfig.RateType)
		if err != nil {
			return err
		}
	}
	return nil
}

func isValidRateType(rateType string) error {
	if rateType != "percentage" && rateType != "fixed" {
		return fmt.Errorf("rate type must be percentage or fixed")
	}
	return nil
}
