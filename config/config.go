package config

import (
	"fmt"
	"sync"

	//"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		JWT  `yaml:"jwt"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"PORT"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `yaml:"pg_url" env:"DATABASE_URL"`
	}

	// JWT -.
	JWT struct {
		Secret     string `yaml:"jwt_secret" env:"JWT_SECRET"`
		EatAuth    int    `yaml:"eatAuth" env:"EAT_AUTH"`
		EatRefresh int    `yaml:"eatRefresh" env:"EAT_REFRESH"`
	}
)

var instance *Config
var once sync.Once

// GetConfig returns app config.
//TODO: Log error
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig("./config/config.yml", instance); err != nil {
			fmt.Println("Config error %w", err)
		}

		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			fmt.Println("Config error %w", err)
		}
	})

	return instance
}
