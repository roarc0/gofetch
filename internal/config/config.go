package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/memory"
)

// Config struct contains the configuration for the application.
type Config struct {
	Collector  collector.CollectorConfig
	Downloader collector.TransmissionConfig
	Memory     memory.Config
	Sources    map[string]collector.Source
	Entries    map[string]filter.Entry
}

func LoadYaml(cfgPath string) (*Config, error) {
	cfgPath = filepath.Join(cfgPath, "config.yaml")

	_, err := os.Stat(cfgPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfgPath = filepath.Join(home, ".config", "gofetch", "config.yaml")
	}

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	cfg.Memory.FilePath = filepath.Join(filepath.Dir(cfgPath), cfg.Memory.FilePath)

	return cfg, nil
}

// Load initializes the configuration.
func Load(cfgPath string) (*Config, error) {
	v, err := LoadViper[Config](cfgPath, "GOFETCH")
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = v.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadViper loads the configuration from the file and environment variables.
func LoadViper[T any](cfgPath, prefix string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	b, err := yaml.Marshal(*new(T))
	if err != nil {
		return nil, err
	}

	defaultConfig := bytes.NewReader(b)
	if err := v.MergeConfig(defaultConfig); err != nil {
		return nil, err
	}

	v.SetConfigName("config")
	v.AddConfigPath(cfgPath)

	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
	}

	v.AutomaticEnv()
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
