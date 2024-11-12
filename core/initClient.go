package core

import (
	"debank_checker_v3/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"log"
	"net/url"
	"time"
)

func GetClient() *fasthttp.Client {
	var dial fasthttp.DialFunc
	randomProxy := utils.GetProxy()

	if randomProxy != "" {
		proxy, err := url.Parse(randomProxy)
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
	}

	return client
}
