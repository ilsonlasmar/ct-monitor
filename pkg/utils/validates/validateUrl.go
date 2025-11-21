package validates

import (
	"net/url"
	"strings"
)

func IsValidURL(rawURL string) bool {
	if strings.Contains(rawURL, "://") {
		_, err := url.Parse(rawURL)
		return err == nil
	}

	prefixedURL := "http://" + rawURL
	u, err := url.Parse(prefixedURL)
	if err != nil {
		return false
	}

	return u.Host != ""
}
