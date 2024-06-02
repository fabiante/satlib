package sdk

import (
	"net/http"
)

type HttpContext struct {
	request  *http.Request
	response http.ResponseWriter
}

func newContext(request *http.Request, response http.ResponseWriter) *HttpContext {
	return &HttpContext{request: request, response: response}
}

func (c *HttpContext) Http() (*Request, *Response) {
	return newRequest(c.request), newResponse(c.response)
}

func (c *HttpContext) HttpStd() (*http.Request, http.ResponseWriter) {
	return c.request, c.response
}
