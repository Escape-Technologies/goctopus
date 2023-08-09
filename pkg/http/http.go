package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	"crypto/sha256"

	"github.com/Escape-Technologies/goctopus/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Post(url string, body []byte) (*Response, error)
	DeleteUrlCache(url string)
}

type client struct {
	// url -> sha256(body) -> response
	cache  map[string]map[string]*Response
	config *config.Config
}

func NewClient(config *config.Config) Client {
	return &client{
		cache:  make(map[string]map[string]*Response),
		config: config,
	}
}

var (
	fastHttpClient *fasthttp.Client
)

func initClient(config *config.Config) {
	fastHttpClient = &fasthttp.Client{
		// MaxConnsPerHost:     1,
		ReadTimeout:         time.Second * time.Duration(config.Timeout),
		WriteTimeout:        time.Second * time.Duration(config.Timeout),
		MaxIdleConnDuration: time.Second * time.Duration(config.Timeout),
		MaxConnDuration:     time.Second * time.Duration(config.Timeout),
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
	if fastHttpClient == nil {
		initClient(c.config)
	}
	sha := sha256.Sum256(body)
	if resp := c.cacheLookup(url, sha); resp != nil {
		return resp, nil
	}
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(body)
	req.SetTimeout(time.Second * time.Duration(c.config.Timeout))
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	log.Debug("Request sent to: ", url)
	err := fastHttpClient.DoTimeout(req, resp, time.Second*time.Duration(c.config.Timeout))
	if err != nil {
		return nil, err
	}
	// The response has to be copied because it will be released after the function returns
	respBody := make([]byte, len(resp.Body()))
	copy(respBody, resp.Body())
	response := &Response{
		StatusCode: resp.StatusCode(),
		Body:       &respBody,
	}
	c.cacheResponse(url, sha, response)
	return response, nil
}

func (c *client) DeleteUrlCache(url string) {
	delete(c.cache, url)
}

func (c *client) cacheLookup(url string, bodySha [32]byte) *Response {
	if _, ok := c.cache[url]; !ok {
		c.cache[url] = make(map[string]*Response)
		return nil
	}
	if resp, ok := c.cache[url][string(bodySha[:])]; ok {
		return resp
	}
	return nil
}

func (c *client) cacheResponse(url string, bodySha [32]byte, resp *Response) {
	if _, ok := c.cache[url]; !ok {
		c.cache[url] = make(map[string]*Response)
	}
	c.cache[url][string(bodySha[:])] = resp
}

func SendToWebhook(url string, body []byte, wg *sync.WaitGroup) error {
	if fastHttpClient == nil {
		log.Panic("Webhook called before any client was initialized")
	}
	defer wg.Done()
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(body)
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	log.Debug("Sending to webhook")
	err := fastHttpClient.Do(req, res)
	if err != nil {
		log.Debugf("error from webhook %v: %v", url, err)
		return err
	}

	if res.StatusCode() != 200 {
		log.Debugf("error from webhook %v: %v", url, res.StatusCode())
		return errors.New("webhook returned non-200 status code")
	}
	return nil
}

func QueryToRequestBody(query string) []byte {
	return []byte(fmt.Sprintf(`{"query":"%s"}`, query))
}
