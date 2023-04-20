package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         uint          `mapstructure:"port"`
	DB           int           `mapstructure:"db"`
	Password     string        `mapstructure:"password"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxPoolSize  int           `mapstructure:"max_pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
}

func NewConfig() *Config {
	initConfig()

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

func initConfig() {
	// set the name of the config file (without extension)
	viper.SetConfigName("default")

	// Set the search paths for the config file
	viper.AddConfigPath("$HOME/go-rate-limiter/deployment/config")

	// Enable support for environment variables
	viper.AutomaticEnv()

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Failed to read config file: %s", err))
	}
}
