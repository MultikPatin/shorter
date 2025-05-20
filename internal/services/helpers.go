package services

import (
	"github.com/google/uuid"
	"main/internal/constants"
	"net/url"
)

var shortPre string

func getKey(u uuid.UUID, p string) string {
	if isURL(p) {
		return u.String()
	}
	return p + u.String()
}

func getResponseLink(k string, p string, h string) string {
	if isURL(p) {
		return p + constants.Delimiter + k + constants.Delimiter
	}
	return h + constants.Delimiter + k + constants.Delimiter
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
