package synthesizer

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	METHOD_GET       string = "GET"
	METHOD_POST      string = "POST"
	SCHEME_HTTP      string = "HTTP"
	SCHEME_HTTPS     string = "HTTPS"
	SCHEME_WEBSOCKET string = "WEBSOCKET"
	PROTO_HTTP_1_1   string = "HTTP/1.1"
	PROTO_HTTP_2     string = "HTTP/2"
)

const (
	KAKAOI_VOICE_STYLE_WOMAN_READ_CALM     string = "WOMAN_READ_CALM"
	KAKAOI_VOICE_STYLE_WOMAN_DIALOG_BRIGHT string = "WOMAN_DIALOG_BRIGHT"
	KAKAOI_VOICE_STYLE_MAN_READ_CALM       string = "MAN_READ_CALM"
	KAKAOI_VOICE_STYLE_MAN_DIALOG_BRIGHT   string = "MAN_DIALOG_BRIGHT"
)

var (
	ErrInvalidScheme error = errors.New("invalid scheme")
	ErrInvalidPath   error = errors.New("invalid path")
	ErrInvalidHost   error = errors.New("invalid host")
	ErrInvalidMethod error = errors.New("invalid method")
)

type RequestApi struct {
	Method string
	Scheme string
	Host   string
	Path   string
}

func (api *RequestApi) GetHttpRequest() (*http.Request, error) {
	var httpRequest *http.Request
	var path string
	if path = api.Method; len(path) == 0 {
		return httpRequest, ErrInvalidPath
	}

	var host string
	if host = api.Host; len(host) == 0 {
		return httpRequest, ErrInvalidHost
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
		return httpRequest, ErrInvalidScheme
	}

	var requestMethod string
	switch strings.ToUpper(api.Method) {
	case METHOD_GET:
		requestMethod = METHOD_GET
	case METHOD_POST:
		requestMethod = METHOD_POST
	default:
		return httpRequest, ErrInvalidMethod
	}

	httpRequest = &http.Request{}
	httpRequest.Method = requestMethod
	httpRequest.URL = &url.URL{Scheme: scheme, Host: host, Path: path}
	httpRequest.Proto = PROTO_HTTP_1_1
	return httpRequest, nil
}
