package aws

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type awsVars struct {
	AWS_ACCESS_KEY_ID     string `mapstructue:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructue:"AWS_SECRET_ACCESS_KEY"`
	AWS_REGION            string `mapstructue:"AWS_REGION"`

	AWS_S3_NAME string `mapstructue:"AWS_S3_NAME"`
	AWS_S3_URI  string `mapstructue:"AWS_S3_URI"`
}

func loadConfig() (config awsVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return awsVars{
			AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AWS_REGION:            os.Getenv("AWS_REGION"),

			AWS_S3_NAME: os.Getenv("AWS_S3_NAME"),
			AWS_S3_URI:  os.Getenv("AWS_S3_URI"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	if config.AWS_S3_URI == "" {
		err = errors.New("AWS_S3_URI is required")
		return
	}

	if config.AWS_ACCESS_KEY_ID == "" {
		err = errors.New("AWS_ACCESS_KEY_ID is required")
		return
	}

	if config.AWS_SECRET_ACCESS_KEY == "" {
		err = errors.New("AWS_SECRET_ACCESS_KEY is required")
		return
	}

	if config.AWS_REGION == "" {
		err = errors.New("AWS_REGION is required")
		return
	}

	if config.AWS_S3_NAME == "" {
		err = errors.New("AWS_S3_NAME is required")
		return
	}

	return
}
