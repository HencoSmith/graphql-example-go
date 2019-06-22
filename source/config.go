package source

import (
	"log"

	"github.com/spf13/viper"

	config "github.com/HencoSmith/graphql-example-go/config/struct"
)

// GetConfig builds and returns a Configuration struct
// Loads from a .yml file
func GetConfig() config.Configuration {
	// Setup config naming & pathing
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	// Attempt to read the specified config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read the config file: %s", err)
	}

	// Config struct to be returned
	var configuration config.Configuration

	// Parse config details into the struct
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Failed to parse configuration file: %v", err)
	}

	return configuration
}
