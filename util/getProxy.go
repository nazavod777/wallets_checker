package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var Proxies []string

// parseProxy normalizes various proxy formats to a valid URL format.
func parseProxy(proxy string) (string, error) {
	patterns := []struct {
		regex    *regexp.Regexp
		template string
	}{
		{regexp.MustCompile(`^([^:@]+):(\d+)$`), "%s://%s:%s"},                               // ip:port
		{regexp.MustCompile(`^((?:http|https|socks4|socks5)://)([^:@]+):(\d+)$`), "%s%s:%s"}, // scheme://ip:port
		{regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):([^:@]+)@([^:@]+):(\d+)$`), "%s://%s:%s@%s:%s"},
		{regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):([^:@]+):([^:@]+):(\d+)$`), "%s://%s:%s@%s:%s"},
		{regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):(\d+)@([^:@]+):([^:@]+)$`), "%s://%s:%s@%s:%s"},
		{regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):(\d+):([^:@]+):([^:@]+)$`), "%s://%s:%s@%s:%s"},
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(proxy)
		if matches == nil {
			continue
		}

		switch len(matches) {
		case 3:
			return fmt.Sprintf(pattern.template, "http", matches[1], matches[2]), nil
		case 4:
			return fmt.Sprintf(pattern.template, matches[1], matches[2], matches[3]), nil
		case 6:
			scheme := strings.TrimSuffix(matches[1], "://")
			if scheme == "" {
				scheme = "http"
			}
			if _, err := strconv.Atoi(matches[3]); err == nil {
				return fmt.Sprintf(pattern.template, scheme, matches[4], matches[5], matches[2], matches[3]), nil
			}
			return fmt.Sprintf(pattern.template, scheme, matches[2], matches[3], matches[4], matches[5]), nil
		}
	}

	return "", fmt.Errorf("invalid proxy format: %s", proxy)
}

// InitProxies loads and normalizes proxies from a file.
func InitProxies(proxyPath string) error {
	lines, err := ReadFileByRows(proxyPath)
	if err != nil {
		return fmt.Errorf("failed to read proxy list: %w", err)
	}

	for _, proxy := range lines {
		parsed, err := parseProxy(proxy)
		if err != nil {
			log.Warnf("failed to parse proxy %q: %v", proxy, err)
			continue
		}
		Proxies = append(Proxies, parsed)
	}

	return nil
}
