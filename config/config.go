package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		User    user
		Metrics metrics
		Swagger swagger
		Tracing tracing
		HTTP    http
		Log     log
		App     app
		JWT     jwt
		Github  github
	}

	github struct {
		ProfileName string `env:"GITHUB_PROFILE_NAME, required"`
		Token       string `env:"GITHUB_TOKEN, required"`
	}

	jwt struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}

	app struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	user struct {
		Password []string `env:"USER_PASSWORD,required"`
	}

	metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	tracing struct {
		Enabled    bool    `env:"TRACING_ENABLED" envDefault:"true"`
		Prettier   bool    `env:"TRACING_PRETTIER" envDefault:"true"`
		SampleRate float64 `env:"TRACING_SAMPLE_RATE" envDefault:"0.1"`
	}

	http struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	log struct {
		Level string `env:"LOG_LEVEL,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("Config Error: %w", err)
	}

	return cfg, nil
}
