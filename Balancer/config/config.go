package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
	Backends   []string         `mapstructure:"backends"`
	DB         DB               `mapstructure:"database"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
}

type HTTPServerConfig struct {
	Port int `mapstructure:"port"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RateLimitConfig struct {
	DefaultCapacity int            `mapstructure:"default_capacity"`
	DefaultRate     float64        `mapstructure:"default_rate"`
	SpecialLimits   []SpecialLimit `mapstructure:"special_limits"`
}

type SpecialLimit struct {
	APIKey   string  `mapstructure:"api_key"`
	Capacity int     `mapstructure:"capacity"`
	Rate     float64 `mapstructure:"rate"`
}

func setDefaults() {
	viper.SetDefault("http_server.port", 8080)
	viper.SetDefault("backends", []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
		"http://localhost:8085",
	})

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "database")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("rate_limit.default_capacity", 100)
	viper.SetDefault("rate_limit.default_rate", 10.0)
	viper.SetDefault("rate_limit.special_limits", []map[string]interface{}{
		{
			"api_key":  "premium_key",
			"capacity": 1000,
			"rate":     100.0,
		},
	})
}

func (c *DB) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
