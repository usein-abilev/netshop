package tools

import "os"

// Gets an environment variable by key if exists, otherwise returns a default value
func TryGetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
