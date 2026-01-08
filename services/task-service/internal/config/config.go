package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

const defaultDBURL = "postgres://postgres:postgres@localhost:5432/go_micro_task_board?sslmode=disable"

// Config holds runtime configuration for the task service.
type Config struct {
	HTTP HTTPConfig `mapstructure:"http"`
	DB   DBConfig   `mapstructure:"db"`
}

type HTTPConfig struct {
	Addr string `mapstructure:"addr"`
}

type DBConfig struct {
	URL string `mapstructure:"url"`
}

// Load reads configuration using defaults, optional config.yaml, and TASKBOARD_ env overrides.
func Load() (Config, error) {
	v := viper.New()
	v.SetDefault("http.addr", ":8082")
	v.SetDefault("db.url", defaultDBURL)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvPrefix("TASKBOARD")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("read config.yaml: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate ensures required fields are present and well-formed.
func (c Config) Validate() error {
	if strings.TrimSpace(c.DB.URL) == "" {
		return fmt.Errorf("db.url is required")
	}
	parsed, err := url.Parse(c.DB.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("db.url must be a valid URL")
	}
	if strings.TrimSpace(c.HTTP.Addr) == "" {
		return fmt.Errorf("http.addr is required")
	}
	return nil
}
