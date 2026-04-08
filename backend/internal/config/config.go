package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "configs/local.yaml"

// Config holds application configuration loaded from YAML with optional env overrides.
type Config struct {
	App      AppConfig      `yaml:"app"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type JWTConfig struct {
	Secret              string `yaml:"secret"`
	AccessTTLSeconds    int    `yaml:"access_ttl_seconds"`
	RefreshTTLSeconds   int    `yaml:"refresh_ttl_seconds"`
}

// Load reads YAML from CONFIG_PATH (default configs/local.yaml), then applies environment overrides.
func Load() (*Config, error) {
	path := strings.TrimSpace(os.Getenv("CONFIG_PATH"))
	if path == "" {
		path = defaultConfigPath
	}

	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read config file %q: %w", path, err)
		}
		cfg = &Config{}
	} else {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	}

	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func applyEnvOverrides(c *Config) {
	if s := strings.TrimSpace(os.Getenv("HTTP_PORT")); s != "" {
		c.App.Port = s
	}
	if s := strings.TrimSpace(os.Getenv("POSTGRES_DSN")); s != "" {
		c.Postgres.DSN = s
	}
	if s := strings.TrimSpace(os.Getenv("REDIS_ADDR")); s != "" {
		c.Redis.Addr = s
	}
	if s := strings.TrimSpace(os.Getenv("JWT_SECRET")); s != "" {
		c.JWT.Secret = s
	}
	if s := strings.TrimSpace(os.Getenv("JWT_ACCESS_TTL_SECONDS")); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			c.JWT.AccessTTLSeconds = n
		}
	}
	if s := strings.TrimSpace(os.Getenv("JWT_REFRESH_TTL_SECONDS")); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			c.JWT.RefreshTTLSeconds = n
		}
	}
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config: nil")
	}
	if strings.TrimSpace(c.App.Port) == "" {
		return errors.New("config: app.port is required (yaml or HTTP_PORT)")
	}
	if strings.TrimSpace(c.Postgres.DSN) == "" {
		return errors.New("config: postgres.dsn is required (yaml or POSTGRES_DSN)")
	}
	if strings.TrimSpace(c.Redis.Addr) == "" {
		return errors.New("config: redis.addr is required (yaml or REDIS_ADDR)")
	}
	if strings.TrimSpace(c.JWT.Secret) == "" {
		return errors.New("config: jwt.secret is required (yaml or JWT_SECRET)")
	}
	if len(strings.TrimSpace(c.JWT.Secret)) < 16 {
		return errors.New("config: jwt.secret must be at least 16 characters")
	}
	if c.JWT.AccessTTLSeconds <= 0 {
		return errors.New("config: jwt.access_ttl_seconds must be positive (yaml or JWT_ACCESS_TTL_SECONDS)")
	}
	if c.JWT.RefreshTTLSeconds <= 0 {
		return errors.New("config: jwt.refresh_ttl_seconds must be positive (yaml or JWT_REFRESH_TTL_SECONDS)")
	}
	if c.JWT.RefreshTTLSeconds < c.JWT.AccessTTLSeconds {
		return errors.New("config: jwt.refresh_ttl_seconds must be >= access_ttl_seconds")
	}
	return nil
}
