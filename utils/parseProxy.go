package utils

import (
	"debank_checker_v3/customTypes"
	"fmt"
	"regexp"
)

func ParseProxy(proxy string) (*customTypes.Proxy, error) {
	re := regexp.MustCompile(`(?i)(?:(socks5|http|https)://)?(?:([a-zA-Z0-9._-]+):([a-zA-Z0-9._-]+)@)?([^:]+):(\d+)`)

	matches := re.FindStringSubmatch(proxy)
	if matches == nil {
		return nil, fmt.Errorf("invalid proxy format")
	}
	proxyData := &customTypes.Proxy{
		Scheme:   matches[1],
		User:     matches[2],
		Password: matches[3],
		IP:       matches[4],
		Port:     matches[5],
	}
	if proxyData.Scheme == "" {
		proxyData.Scheme = "http"
	}

	return proxyData, nil
}
