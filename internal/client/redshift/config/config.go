package config

import (
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	Port               int
	Password           string
	User               string
	Host               string
	Name               string
	ServiceName        string
	MaxConns           int
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

func NewFromEnv() (*Config, error) {
	cfg := &Config{
		Port:               env.GetInt("REDSHIFT_PORT", 5439),
		Password:           env.GetString("REDSHIFT_PASSWORD", ""),
		User:               env.GetString("REDSHIFT_USER", ""),
		Host:               env.GetString("REDSHIFT_HOST", ""),
		Name:               env.GetString("REDSHIFT_NAME", ""),
		ServiceName:        env.GetString("SERVICE_NAME", ""),
		MaxConns:           env.GetInt("REDSHIFT_MAX_CONNECTIONS", 20),
		AwsAccessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
		AwsSecretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Password, validation.Required, validation.Required.Error("REDSHIFT_PASSWORD was not specified")),
		validation.Field(&c.User, validation.Required, validation.Required.Error("REDSHIFT_USER was not specified")),
		validation.Field(&c.Host, validation.Required, validation.Required.Error("REDSHIFT_HOST was not specified")),
		validation.Field(&c.Name, validation.Required, validation.Required.Error("REDSHIFT_NAME was not specified")),
	)
}
