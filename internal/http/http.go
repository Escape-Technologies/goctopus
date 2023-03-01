package http

import (
	"crypto/tls"
	"time"

	"github.com/Escape-Technologies/goctopus/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Post(url string, body []byte) ([]byte, error)
}

type client struct{}

func NewClient() Client {
	return &client{}
}

var (
	fastHttpClient *fasthttp.Client
)

func init() {
	fastHttpClient = &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func (c *client) Post(url string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(body)
	req.SetTimeout(time.Second * time.Duration(config.Conf.Timeout))
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	log.Debug("Request sent to: ", url)
	err := fastHttpClient.DoTimeout(req, resp, time.Second*time.Duration(config.Conf.Timeout))
	if err != nil {
		log.Debugf("Error from %v: %v", url, err)
		return nil, err
	}
	log.Debugf("Response from %v: %v", url, string(resp.Body()))
	return resp.Body(), nil
}
