package sdk

import (
	"net/http"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
}

func newContext(request *http.Request, response http.ResponseWriter) *Context {
	return &Context{request: request, response: response}
}

func (c *Context) Http() (*Request, *Response) {
	return newRequest(c.request), newResponse(c.response)
}

func (c *Context) HttpStd() (*http.Request, http.ResponseWriter) {
	return c.request, c.response
}
