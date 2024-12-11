package core

import (
	"crypto/tls"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"log"
	"net/url"
	"time"
)

func GetClient(proxy string) *fasthttp.Client {
	var dial fasthttp.DialFunc

	if proxy != "" {
		proxy, err := url.Parse(proxy)
		if err != nil {
			log.Panicf("Error Unparsing Proxy: %v\n", err)
		}

		switch proxy.Scheme {
		case "http", "https":
			dial = fasthttpproxy.FasthttpHTTPDialer(proxy.String())
		case "socks4":
			dial = fasthttpproxy.FasthttpSocksDialer(proxy.String())
		case "socks5":
			dial = fasthttpproxy.FasthttpSocksDialer(proxy.String())
		default:
			log.Panicf("Unsupported proxy scheme: %s\n", proxy.Scheme)
		}
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},

		Renegotiation:          tls.RenegotiateNever,
		SessionTicketsDisabled: false,
		InsecureSkipVerify:     true,
	}

	client := &fasthttp.Client{
		Dial:                          dial,
		MaxConnsPerHost:               0,
		MaxIdleConnDuration:           90 * time.Second,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		ReadTimeout:                   15 * time.Second,
		WriteTimeout:                  15 * time.Second,
		MaxConnWaitTimeout:            15 * time.Second,
		StreamResponseBody:            true,
		TLSConfig:                     tlsConfig,
	}

	return client
}
