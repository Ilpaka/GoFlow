package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type PostgresConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	TokenTTLSeconds int    `mapstructure:"token_ttl_seconds"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	applyEnvOverrides(&c)

	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

func applyEnvOverrides(c *Config) {
	if s := os.Getenv("HTTP_PORT"); s != "" {
		c.App.Port = strings.TrimSpace(s)
	}
	if s := os.Getenv("POSTGRES_DSN"); s != "" {
		c.Postgres.DSN = strings.TrimSpace(s)
	}
	if s := os.Getenv("REDIS_ADDR"); s != "" {
		c.Redis.Addr = strings.TrimSpace(s)
	}
	if s, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		c.Redis.Password = s
	}
	if s := os.Getenv("REDIS_DB"); s != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			c.Redis.DB = n
		}
	}
	if s := os.Getenv("JWT_SECRET"); s != "" {
		c.JWT.Secret = strings.TrimSpace(s)
	}
	if s := os.Getenv("JWT_TOKEN_TTL_SECONDS"); s != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			c.JWT.TokenTTLSeconds = n
		}
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.App.Port) == "" {
		return errors.New("config: app.port is required")
	}
	if strings.TrimSpace(c.Postgres.DSN) == "" {
		return errors.New("config: postgres.dsn is required")
	}
	if strings.TrimSpace(c.Redis.Addr) == "" {
		return errors.New("config: redis.addr is required")
	}
	if strings.TrimSpace(c.JWT.Secret) == "" {
		return errors.New("config: jwt.secret is required")
	}
	if len(strings.TrimSpace(c.JWT.Secret)) < 16 {
		return errors.New("config: jwt.secret must be at least 16 characters")
	}
	if c.JWT.TokenTTLSeconds <= 0 {
		return errors.New("config: jwt.token_ttl_seconds must be positive")
	}
	return nil
}
