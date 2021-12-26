package synthesizer

import (
	"net/http"
)

// <-- KAKAOI REQUEST API -->
const (
	KAKAOI_METHOD       string = http.MethodPost
	KAKAOI_SCHEME       string = "https"
	KAKAOI_HOST         string = "kakaoi-newtone-openapi.kakao.com"
	KAKAOI_PATH         string = "/v1/synthesize"
	KAKAOI_CONTENT_TYPE string = "application/xml"
)

type KakaoTTSRequest struct {
	RequestApi
}
