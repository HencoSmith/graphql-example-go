package config

import (
	"time"
)

// JWTConfiguration relates to JWT variables
type JWTConfiguration struct {
	Key        string
	Expiration time.Duration
}
