package cmdln

import (
	"os"
	"strings"
)

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

// ContainsIgnoreCase checks if string `b` is a substring of string `a`, ignoring case.
//
// Parameters:
//
//	a - The string to search within.
//	b - The substring to search for.
//
// Returns:
//
//	bool - true if `b` is a case-insensitive substring of `a`, otherwise false.
func ContainsIgnoreCase(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}
