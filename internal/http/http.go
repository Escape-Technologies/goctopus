package http

import (
	"crypto/tls"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	Client *fasthttp.Client
)

func init() {
	Client = &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func Post(url string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(body)
	req.SetTimeout(time.Second * 2)
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	log.Debug("Request sent to: ", url)
	err := Client.DoTimeout(req, resp, time.Second*2)
	if err != nil {
		log.Debugf("Error from %v: %v", url, err)
		return nil, err
	}
	log.Debugf("Response from %v: %v", url, resp.Body())
	return resp.Body(), nil
}
