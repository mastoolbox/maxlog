package cmdln

import "os"

const (
	// MAS Manage
	DefaultLabels = "all, ui, cron, mea, rpt, jms"
	AppTypeName   = "mas.ibm.com/appTypeName"
)

// GetEnv retrieves the value of an environment variable, returning a default value if the variable is not set.
//
// Parameters:
//
//	key          - The name of the environment variable to retrieve.
//	defaultValue - The default value to return if the environment variable is not set or is empty.
//
// Returns:
//
//	string - The value of the environment variable if set, otherwise the default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
