package utils

import (
	"net/url"
)

func GetPathFromURI(uri string) (string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}
