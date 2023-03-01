package helpers

import "github.com/valyala/fasthttp"

func MockHttpResponse(statusCode int, body string) *fasthttp.Response {
	resp := &fasthttp.Response{}
	resp.SetStatusCode(statusCode)
	resp.AppendBody([]byte(body))
	return resp
}
