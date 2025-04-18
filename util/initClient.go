package util

import (
	"debank_checker_v3/global"
	"net/url"
	"sync/atomic"

	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	log "github.com/sirupsen/logrus"
)

var currentClientIndex uint64

// CreateClient initializes a tls-client with the given proxy (optional).
func CreateClient(proxy string) tlsclient.HttpClient {
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(10),
		tlsclient.WithClientProfile(profiles.Chrome_124),
		tlsclient.WithNotFollowRedirects(),
		tlsclient.WithInsecureSkipVerify(),
	}

	if proxy != "" {
		parsedURL, err := url.Parse(proxy)
		if err != nil {
			log.Panicf("failed to parse proxy %q: %v", proxy, err)
		}
		options = append(options, tlsclient.WithProxyUrl(parsedURL.String()))
	}

	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		log.Panicf("failed to create TLS client: %v", err)
	}

	return client
}

// GetClient returns the next client from global.Clients in round-robin order.
// It is thread-safe and cycles through all clients.
func GetClient() tlsclient.HttpClient {
	total := uint64(len(global.Clients))
	if total == 0 {
		log.Panic("no clients available in global.Clients")
	}

	index := atomic.AddUint64(&currentClientIndex, 1)
	return global.Clients[index%total]
}
