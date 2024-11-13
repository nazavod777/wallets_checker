package utils

import (
	"fmt"
	"regexp"
)

var proxyFormatsRegexp = []*regexp.Regexp{
	// protocol://login:password@host:port
	regexp.MustCompile(`^(?:(?P<protocol>.+)://)?(?P<login>[^:]+):(?P<password>[^@|:]+)[@|:](?P<host>[^:]+):(?P<port>\d+)$`),
	// protocol://host:port@login:password
	regexp.MustCompile(`^(?:(?P<protocol>.+)://)?(?P<host>[^:]+):(?P<port>\d+)[@|:](?P<login>[^:]+):(?P<password>[^:]+)$`),
	// host:port
	regexp.MustCompile(`^(?:(?P<protocol>.+)://)?(?P<host>[^:]+):(?P<port>\d+)$`),
}

func ParseProxy(proxy string) (string, error) {
	for _, pattern := range proxyFormatsRegexp {
		match := pattern.FindStringSubmatch(proxy)
		if match != nil {
			protocol := match[1]
			if protocol == "" {
				protocol = "http"
			}
			login := match[2]
			password := match[3]
			host := match[4]
			port := match[5]

			if login != "" && password != "" {
				return fmt.Sprintf("%s://%s:%s@%s:%s", protocol, login, password, host, port), nil
			}
			return fmt.Sprintf("%s://%s:%s", protocol, host, port), nil
		}
	}

	return "", fmt.Errorf("invalid proxy format")
}
