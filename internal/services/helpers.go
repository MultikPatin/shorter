package services

import (
	"github.com/google/uuid"
	"net/url"
)

const (
	delimiter = "/"
)

func GetDBKey(u uuid.UUID, p string) string {
	if IsURL(p) {
		return u.String()
	}
	return p + u.String()
}

func GetResponseLink(k string, p string, h string) string {
	if IsURL(p) {
		return p + delimiter + k + delimiter
	}
	return h + delimiter + k + delimiter
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
