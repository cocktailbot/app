package config

import (
	"os"
)

// Get particular config value
func Get(name string) string {
	return os.Getenv(name)
}
