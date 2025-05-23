package services // Package services provides helper functions for generating keys and URLs.

import (
	"github.com/google/uuid"
	"main/internal/constants"
	"net/url"
)

// shortPre represents a configurable prefix for generated short links.
var shortPre string

// getKey generates a unique key for a given UUID and prefix.
// If the prefix is a valid URL, the key includes only the UUID.
// Otherwise, the key combines the prefix and UUID.
func getKey(u uuid.UUID, p string) string {
	if isURL(p) {
		return u.String()
	}
	return p + u.String()
}

// getResponseLink constructs a full response URL combining the key, prefix, and host.
// Behavior depends on whether the prefix is a valid URL.
func getResponseLink(k string, p string, h string) string {
	if isURL(p) {
		return p + constants.Delimiter + k + constants.Delimiter
	}
	return h + constants.Delimiter + k + constants.Delimiter
}

// isURL determines if a given string is a well-formed URL.
func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
