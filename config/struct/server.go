package config

import (
	"time"
)

// ServerConfiguration relates to server variables
type ServerConfiguration struct {
	Port    string
	Timeout time.Duration
}
