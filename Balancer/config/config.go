package config

type Config struct {
	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
	Backends   []string         `mapstructure:"backends"`
}

type HTTPServerConfig struct {
	Port int `mapstructure:"port"`
}
