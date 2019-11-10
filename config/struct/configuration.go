package config

// Configuration Links all sub configurations"
type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	JWT      JWTConfiguration
}
