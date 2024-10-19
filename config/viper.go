package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	PORT string `mapstructure:"PORT"`
	MONGODB_URI string `mapstructure:"MONGODB_URI"`
	MONGODB_NAME string `mapstructure:"MONGODB_NAME"`
	AUTH0_DOMAIN string `mapstructure:"AUTH0_DOMAIN"`
	AUTH0_AUDIENCE string `mapstructure:"AUTH0_AUDIENCE"`
}

func LoadConfig() (config EnvVars, err error){
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			PORT: os.Getenv("PORT"),
			MONGODB_URI: os.Getenv("MONGODB_URI"),
			MONGODB_NAME: os.Getenv("MONGODB_NAME"),
			AUTH0_DOMAIN: os.Getenv("AUTH0_DOMAIN"),
			AUTH0_AUDIENCE: os.Getenv("AUTH0_AUDIENCE"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if (err != nil) {
		return
	}

	err = viper.Unmarshal(&config)

	if config.PORT == "" {
		err = errors.New("PORT is required")
		return
	}

	if config.MONGODB_URI == "" {
		err = errors.New("MONGODB_URI is required")
		return
	}

	if config.MONGODB_NAME == "" {
		err = errors.New("MONGODB_NAME is required")
		return
	}

	if config.AUTH0_DOMAIN == "" {
		err = errors.New("AUTH0_DOMAIN is required")
		return
	}

	if config.AUTH0_AUDIENCE == "" {
		err = errors.New("AUTH0_AUDIENCE is required")
		return
	}

	return
}