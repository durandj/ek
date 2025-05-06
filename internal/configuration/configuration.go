package configuration

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Environment Environment `default:"production"`

	Logging LoggingConfiguration
	HTTP    HTTPConfiguration
	Source  SourceConfiguration
}

func NewConfigurationFromEnv() (Configuration, error) {
	var config Configuration

	if err := envconfig.Process("ek", &config); err != nil {
		return Configuration{}, fmt.Errorf(
			"unable to load configuration from environment: %w",
			err,
		)
	}

	return config, nil
}

type Environment string

const (
	EnvironmentInvalid     Environment = ""
	EnvironmentProduction  Environment = "production"
	EnvironmentDevelopment Environment = "development"
)

func (env Environment) String() string {
	return string(env)
}
