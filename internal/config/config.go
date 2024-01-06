package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Addr     string        `mapstructure:"ADDRESS"`
	Duration time.Duration `mapstructure:"DURATION"`
	Secret   string        `mapstructure:"SECRET"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	return config, nil
}
