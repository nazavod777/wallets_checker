package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseProxy(proxy string) (string, error) {
	patterns := []struct {
		regex    *regexp.Regexp
		template string
	}{
		// ip:port
		{
			regexp.MustCompile(`^([^:@]+):(\d+)$`),
			"%s://%s:%s",
		},
		// scheme://ip:port
		{
			regexp.MustCompile(`^((?:http|https|socks4|socks5)://)([^:@]+):(\d+)$`),
			"%s%s:%s",
		},
		// scheme://user:pass@ip:port
		{
			regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):([^:@]+)@([^:@]+):(\d+)$`),
			"%s://%s:%s@%s:%s",
		},
		// scheme://user:pass:ip:port
		{
			regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):([^:@]+):([^:@]+):(\d+)$`),
			"%s://%s:%s@%s:%s",
		},
		// scheme://ip:port@user:pass
		{
			regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):(\d+)@([^:@]+):([^:@]+)$`),
			"%s://%s:%s@%s:%s",
		},
		// scheme://ip:port:user:pass
		{
			regexp.MustCompile(`^((?:http|https|socks4|socks5)://)?([^:@]+):(\d+):([^:@]+):([^:@]+)$`),
			"%s://%s:%s@%s:%s",
		},
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(proxy)
		if matches == nil {
			continue
		}

		switch len(matches) {
		case 3: // Простой формат ip:port
			return fmt.Sprintf(pattern.template, "http", matches[1], matches[2]), nil
		case 4: // Формат scheme://ip:port
			return fmt.Sprintf(pattern.template, matches[1], matches[2], matches[3]), nil
		case 6: // Форматы с user:pass
			scheme := matches[1]
			if scheme == "" {
				scheme = "http://"
			}
			scheme = strings.TrimSuffix(scheme, "://")

			if strings.Contains(pattern.template, "@") {
				if isPort(matches[3]) {
					return fmt.Sprintf(pattern.template, scheme, matches[4], matches[5], matches[2], matches[3]), nil
				}
				return fmt.Sprintf(pattern.template, scheme, matches[2], matches[3], matches[4], matches[5]), nil
			}
		}
	}

	return "", fmt.Errorf("invalid proxy format: %s", proxy)
}

func isPort(s string) bool {
	match, _ := regexp.MatchString(`^\d+$`, s)
	return match
}
