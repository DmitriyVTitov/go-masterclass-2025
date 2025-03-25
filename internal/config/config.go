package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel          int    `mapstructure:"log_level"`
	TelemetryEndpoint string `mapstructure:"telemetry_endpoint"`
	AppBaseURL        string `mapstructure:"app_base_url"`
	AppBasePath       string `mapstructure:"app_base_path"`
	Server            Server `mapstructure:"server"`
}

type Server struct {
	ListenPort int `mapstructure:"listen_port"`
}

func Load(path string) (*Config, error) {
	setDefaults()

	var c Config

	if path != "" {
		_, err := os.Lstat(path)
		if err != nil {
			return nil, err
		}
		viper.SetConfigFile(path)

		err = viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	envVars(&c)

	return &c, nil
}

// envVars sets config values from env vars.
func envVars(c *Config) {
	viper.AutomaticEnv()

	// Log
	if val := viper.GetString("LOG_LEVEL"); val != "" {
		if intVal := viper.GetInt("LOG_LEVEL"); intVal >= 0 {
			c.LogLevel = intVal
		}
	}

	// Jaeger
	if val := viper.GetString("TELEMETRY_ENDPOINT"); val != "" {
		c.TelemetryEndpoint = val
	}

	// Swagger
	if val := viper.GetString("APP_BASE_URL"); val != "" {
		c.AppBaseURL = val
	}
	if val := viper.GetString("APP_BASE_PATH"); val != "" {
		c.AppBasePath = val
	}
}

func setDefaults() {
	viper.SetDefault("log_level", 1)
	viper.SetDefault("server.listen_port", 8080)
}
