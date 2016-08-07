package routing

import (
  "net/http"
)

type Context interface {
  Request() *http.Request
  SetResponseData(data interface{})
  SetResponseHeader(string, string)
  ResponseData() interface{}
  ResponseHeaders() map[string]string
  Status() int
}

type BasicContext struct {
  Context
  request *http.Request
  response http.ResponseWriter
  responseHeaders map[string]string
  status int
  stop bool
  data interface{}
}

func (context *BasicContext) Request() *http.Request {
  return context.request
}

func (context *BasicContext) SetResponseData(data interface{}) {
  context.data = data
}

func (context *BasicContext) SetResponseHeader(name string, value string) {
  context.responseHeaders[name] = value
}

func (context *BasicContext) ResponseData() interface{} {
  return context.data
}

func (context *BasicContext) ResponseHeaders() map[string]string {
  return context.responseHeaders
}

func (context *BasicContext) Status() int {
  if context.status == 0 {
    return http.StatusOK
  } else {
    return context.status
  }
}
