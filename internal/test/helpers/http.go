package helpers

import "github.com/Escape-Technologies/goctopus/internal/http"

func MockHttpResponse(statusCode int, body string) *http.Response {
	bodyBytes := []byte(body)
	resp := &http.Response{
		Body:       &bodyBytes,
		StatusCode: statusCode,
	}
	return resp
}
