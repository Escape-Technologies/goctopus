package http

import (
	"math/rand"
	"time"

	"net"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func RotatingProxyDialerTimeout(proxies []string, timeout time.Duration) fasthttp.DialFunc {
	if len(proxies) == 0 {
		// use the default dialer
		log.Infof("No proxies set, using default dialer")
		return nil
	}

	log.Infof("Using %d proxies", len(proxies))

	dialFunctions := make([]fasthttp.DialFunc, len(proxies))
	for i, proxy := range proxies {
		// @todo we could reimplement their FasthttpHTTPDialerTimeout to keep the connections alive with the proxies
		dialFunctions[i] = fasthttpproxy.FasthttpHTTPDialerTimeout(proxy, timeout)
	}

	return func(addr string) (net.Conn, error) {
		return dialFunctions[rand.Intn(len(dialFunctions))](addr)
	}
}
