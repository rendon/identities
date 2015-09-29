package config

import (
	"time"
)

var MongoDSN = `mongodb-server`

var ExpirationTime = map[string]time.Duration{
	"twitter":   30 * 24 * time.Hour,
	"instagram": 30 * 24 * time.Hour,
}

var Networks = map[string]bool{
	"twitter":   true,
	"instagram": true,
}
