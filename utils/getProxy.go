package utils

import (
	"fmt"
	"log"
	"math/rand"
)

var Proxies []string

func InitProxies() error {
	proxiesFile, err := ReadFileByRows("./data/proxies.txt")

	if err != nil {
		return fmt.Errorf("error When Reading Proxy: %s", err)
	}

	for _, proxy := range proxiesFile {
		parsedProxy, err := ParseProxy(proxy)

		if err != nil {
			log.Printf("Error When Parsing Proxy %s: %s", proxy, err)
			continue
		}

		Proxies = append(Proxies, parsedProxy)
	}

	return nil
}

func GetProxy() string {
	if len(Proxies) == 0 {
		return ""
	}
	return Proxies[rand.Intn(len(Proxies))]
}
