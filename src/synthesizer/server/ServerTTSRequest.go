package synthesizer

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	METHOD_GET       = "GET"
	METHOD_POST      = "POST"
	SCHEME_HTTP      = "HTTP"
	SCHEME_HTTPS     = "HTTPS"
	SCHEME_WEBSOCKET = "WEBSOCKET"
	PROTO_HTTP_1_1   = "HTTP/1.1"
	PROTO_HTTP_2     = "HTTP/2"
)

var ERROR_INVALID_SCHEME error = errors.New("invalid scheme")
var ERROR_INVALID_PATH error = errors.New("invalid path")
var ERROR_INVALID_HOST error = errors.New("invalid host")
var ERROR_INVALID_METHOD error = errors.New("invalid method")

type RequestApi struct {
	Method string
	Scheme string
	Host   string
	Path   string
}

type ServerTTSRequest interface {
	NewHttpRequest(message string) (*http.Request, error)
}

func (api *RequestApi) NewHttpRequest() (*http.Request, error) {
	var httpRequest *http.Request
	var path string
	if path = api.Method; len(path) == 0 {
		return httpRequest, ERROR_INVALID_PATH
	}

	var host string
	if host = api.Host; len(host) == 0 {
		return httpRequest, ERROR_INVALID_HOST
	}

	var scheme string
	switch strings.ToUpper(api.Scheme) {
	case SCHEME_HTTP:
		scheme = strings.ToLower(SCHEME_HTTP)
	case SCHEME_HTTPS:
		scheme = strings.ToLower(SCHEME_HTTPS)
	case SCHEME_WEBSOCKET:
		scheme = strings.ToLower(SCHEME_WEBSOCKET)
	default:
		return httpRequest, ERROR_INVALID_SCHEME
	}

	var requestMethod string
	switch strings.ToUpper(api.Method) {
	case METHOD_GET:
		requestMethod = METHOD_GET
	case METHOD_POST:
		requestMethod = METHOD_POST
	default:
		return httpRequest, ERROR_INVALID_METHOD
	}

	httpRequest = &http.Request{}
	httpRequest.Method = requestMethod
	httpRequest.URL = &url.URL{Scheme: scheme, Host: host, Path: path}
	httpRequest.Proto = PROTO_HTTP_1_1
	return httpRequest, nil
}
