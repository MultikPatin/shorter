package services

import (
	"github.com/google/uuid"
	"net/url"
)

const (
	urlPrefix = "http://"
	delimiter = "/"
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
		return p + delimiter + k + delimiter
	}
	return h + delimiter + k + delimiter
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
