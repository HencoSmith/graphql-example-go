package config

// DatabaseConfiguration relates to server variables
type DatabaseConfiguration struct {
	User     string
	Host     string
	Port     string
	Name     string
	Password string
	SSL      string
}
