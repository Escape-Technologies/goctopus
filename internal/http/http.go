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
		MaxConnsPerHost: 1,
		// MaxIdleConnDuration: time.Second * time.Duration(config.Conf.Timeout),
		// MaxConnDuration:     time.Second * time.Duration(config.Conf.Timeout),

		// Might need to implement this
		// MaxResponseBodySize: 1024 * 256,
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
		return nil, err
	}
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
	if config.Conf.WebhookUrl == "" {
		return nil
	}
	AvoidNetworkCongestion()
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

// func AvoidNetworkCongestion() {
// 	var (
// 		url          = "http://8.8.8.8"
// 		maxBackoff   = 10
// 		backoffCount = 0
// 		backoffTime  = 500 * time.Millisecond
// 	)
// 	for {
// 		req, err := http.NewRequest("HEAD", url, nil)
// 		if err != nil {
// 			panic(err)
// 		}
// 		resp, err := http.DefaultTransport.RoundTrip(req)
// 		log.Infof("Network congestion avoidance: %v", resp.StatusCode)
// 		if err != nil {
// 			backoff(&backoffCount, maxBackoff, backoffTime)
// 		}
// 		if resp.StatusCode == 200 {
// 			break
// 		}
// 		backoff(&backoffCount, maxBackoff, backoffTime)
// 	}
// }

// Dialing doest seem to work
// func AvoidNetworkCongestion() {
// 	var (
// 		addr         = "google.com:80"
// 		maxBackoff   = 10
// 		backoffCount = 0
// 		backoffTime  = 500 * time.Millisecond
// 	)
// 	for {
// 		con, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(config.Conf.Timeout))
// 		if err != nil {
// 			backoff(&backoffCount, maxBackoff, backoffTime)
// 		}
// 		if con != nil {
// 			con.Close()
// 			break
// 		}
// 	}
// }

func AvoidNetworkCongestion() {}

func backoff(backoffCount *int, maxBackoff int, backoffTime time.Duration) {
	log.Warnf("Backoff %v", *backoffCount)
	if *backoffCount == maxBackoff {
		panic("Max backoff reached")
	}
	*backoffCount++
	time.Sleep(backoffTime * time.Duration(*backoffCount))
}
