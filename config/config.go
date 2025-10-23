package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Env         string `mapstructure:"ENV"`
	DSN         string `mapstructure:"DSN"`
	PORT        string `mapstructure:"PORT"`
	GCPProject  string `mapstructure:"GCP_PROJECT"`
	GCPRegion   string `mapstructure:"GCP_REGION"`

	ClerkAPIKey string `mapstructure:"CLERK_API_KEY"`
	ConsulHost  string `mapstructure:"CONSUL_HOST"`
	ConsulPort  string `mapstructure:"CONSUL_PORT"`

	HouseholdHost string `mapstructure:"HOUSEHOLD_SERVICE_HOST"`
	MemberHost    string `mapstructure:"MEMBER_SERVICE_HOST"`
	MealHost      string `mapstructure:"MEAL_SERVICE_HOST"`
	EventHost     string `mapstructure:"EVENT_SERVICE_HOST"`
}

func LoadConfig(path string) (Config, error) {
	// Set defaults first
	viper.SetDefault("CLERK_API_KEY", "")
	viper.SetDefault("CONSUL_HOST", "consul")
	viper.SetDefault("CONSUL_PORT", "8500")
	viper.SetDefault("DSN", "postgresql://postgres:postgres@127.0.0.1:54322/postgres")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SERVICE_NAME", "gateway")
	viper.SetDefault("HO", "http://localhost:8080")
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("GCP_PROJECT", "")
	viper.SetDefault("GCP_REGION", "")

	// Configure viper
	viper.AddConfigPath(path)
	viper.SetConfigName(".env.local")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file found but error reading it
			return Config{}, err
		}
		// Config file not found - continue with env vars and defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
