package http

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/Escape-Technologies/goctopus/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Post(url string, body []byte) (*Response, error)
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

type Response struct {
	StatusCode int
	Body       *[]byte
}

func (c *client) Post(url string, body []byte) (*Response, error) {
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
	log.Debugf("Recived status code %v: %v", resp.StatusCode(), string(resp.Body()))
	// The response has to be copied because it will be released after the function returns
	respBody := make([]byte, len(resp.Body()))
	copy(respBody, resp.Body())
	response := &Response{
		StatusCode: resp.StatusCode(),
		Body:       &respBody,
	}
	return response, nil
}

func SendToWebhook(body []byte, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Debugf("webhook url: %v", config.Conf.WebhookUrl)

	if config.Conf.WebhookUrl == "" {
		return nil
	}
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(config.Conf.WebhookUrl)
	req.SetBody(body)
	defer fasthttp.ReleaseRequest(req)
	log.Debug("Sending to webhook")
	err := fastHttpClient.Do(req, nil)
	if err != nil {
		log.Debugf("Error from %v: %v", config.Conf.WebhookUrl, err)
	}
	return err
}
