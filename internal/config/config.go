package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var config struct {
	App struct {
		Port int `env:"PORT" env-default:"8000" env-description:"HTTP port of the API service"`
	} `env-prefix:"APP_"`
	Database struct {
		Host     string `env:"HOST" env-default:"localhost" env-description:"main database hostname"`
		Port     uint16 `env:"PORT" env-default:"5432" env-description:"main database port"`
		Name     string `env:"NAME" env-default:"provisioning" env-description:"main database name"`
		User     string `env:"USER" env-default:"postgres" env-description:"main database username"`
		Password string `env:"PASSWORD" env-default:"" env-description:"main database password"`
	} `env-prefix:"DATABASE_"`
	Logging struct {
		Level         string `env:"LEVEL" env-default:"info" env-description:"logger level (trace, debug, info, warn, error, fatal, panic)"`
		DatabaseLevel string `env:"DB_LEVEL" env-default:"info" env-description:"database logs level (trace, debug, info, warn, error, fatal, panic)"`
	} `env-prefix:"LOGGING_"`
	Cloudwatch struct {
		Enabled bool   `env:"ENABLED" env-default:"false" env-description:"cloudwatch logging exporter (enabled in clowder)"`
		Region  string `env:"REGION" env-default:"" env-description:"cloudwatch logging AWS region"`
		Key     string `env:"KEY" env-default:"" env-description:"cloudwatch logging key"`
		Secret  string `env:"SECRET" env-default:"" env-description:"cloudwatch logging secret"`
		Session string `env:"SESSION" env-default:"" env-description:"cloudwatch logging session"`
		Group   string `env:"GROUP" env-default:"" env-description:"cloudwatch logging group"`
	} `env-prefix:"CLOUDWATCH_"`
}

var (
	Application = &config.App
	Database    = &config.Database
	Logging     = &config.Logging
	Cloudwatch  = &config.Cloudwatch
)

// Initialize loads configuration from provided .env files.
// Files are applied one by one.
// If one config is defined by multiple files, last file wins.
// When no file out of provided exists, the env variables are loaded.
func Initialize(configFiles ...string) {
	var loaded bool
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			// if config file exists, load it (also loads environmental variables)
			if readErr := cleanenv.ReadConfig(configFile, &config); readErr != nil {
				panic(readErr)
			}
			loaded = true
		}
	}

	// if no file found use environmental variables instead
	if !loaded {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			panic(err)
		}
	}

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		// database
		config.Database.Host = cfg.Database.Hostname
		config.Database.Port = uint16(cfg.Database.Port)
		config.Database.User = cfg.Database.Username
		config.Database.Password = cfg.Database.Password
		config.Database.Name = cfg.Database.Name

		// cloudwatch (is blank in ephemeral)
		cw := cfg.Logging.Cloudwatch
		if cw.Region != "" && cw.AccessKeyId != "" && cw.SecretAccessKey != "" && cw.LogGroup != "" {
			config.Cloudwatch.Enabled = true
			config.Cloudwatch.Key = cw.AccessKeyId
			config.Cloudwatch.Secret = cw.SecretAccessKey
			config.Cloudwatch.Region = cw.Region
			config.Cloudwatch.Group = cw.LogGroup
		}
	}
}

func HelpText() (string, error) {
	headerText := ""
	text, err := cleanenv.GetDescription(&config, &headerText)
	if err != nil {
		return "", fmt.Errorf("cannot generate help text: %w", err)
	}
	return text, nil
}
