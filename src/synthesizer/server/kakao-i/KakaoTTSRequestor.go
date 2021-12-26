package synthesizer

import (
	"errors"
	"io"
	"net/http"

	synthesizer "github.com/PaulianBio/src/synthesizer/server"
)

var ERROR_INVALID_CONTENT_TYPE error = errors.New("invalid Content-Type")
var ERROR_INVALID_AUTHORIZATION error = errors.New("invalid Authorization")

/*
https://developers.kakao.com/docs/latest/ko/voice/rest-api#text-to-speech
* POST /v1/recognize HTTP/1.1
* Host: kakaoi-newtone-openapi.kakao.com
* Transfer-Encoding: chunked
* Content-Type: application/octet-stream
* Authorization: KakaoAK {REST_API_KEY}
*/

type KakaoTTSRequestor struct {
	KakaoTTSRequest
	contentType   string
	authorization string
}

func CreateKakaoTTSRequestor(authKey string) *KakaoTTSRequestor {
	var requestor *KakaoTTSRequestor
	requestor = new(KakaoTTSRequestor)
	requestor.KakaoTTSRequest.RequestApi = synthesizer.RequestApi{
		Method: KAKAOI_METHOD,
		Scheme: KAKAOI_SCHEME,
		Host:   KAKAOI_HOST,
		Path:   KAKAOI_PATH,
	}
	requestor.contentType = KAKAOI_CONTENT_TYPE
	requestor.authorization = authKey
	return requestor
}

func (requestor *KakaoTTSRequestor) RequestSynthesizer(messages map[string][]string) (*http.Response, error) {
	var httpRequest *http.Request
	var err error

	// httpRequest, err = requestor.KakaoTTSRequest.NewHttpRequest(messages)
	// if err != nil {
	// 	return nil, err
	// }

	var authorization string
	if len(requestor.authorization) == 0 {
		return nil, ERROR_INVALID_CONTENT_TYPE
	}
	authorization = "KakaoAK" + " " + requestor.authorization

	var contentType string
	if len(requestor.contentType) == 0 {
		return nil, ERROR_INVALID_AUTHORIZATION
	}
	contentType = requestor.contentType

	httpRequest.Header = map[string][]string{
		"Authorization": {authorization},
		"Content-Type":  {contentType},
	}

	client := http.DefaultClient
	var httpResponse *http.Response
	if httpResponse, err = client.Do(httpRequest); err != nil {
		return nil, err
	}

	return httpResponse, nil
}

type KakaoVoiceMessage struct {
	Style   string
	Message string
}

func (requestor *KakaoTTSRequestor) RequestSynthesizerAsync(messages []KakaoVoiceMessage, pipeWriter *io.PipeWriter) error {
	var httpRequest *http.Request
	var err error

	httpRequest, err = requestor.KakaoTTSRequest.NewHttpRequest(&messages)
	if err != nil {
		return err
	}

	var authorization string
	if len(requestor.authorization) == 0 {
		return ERROR_INVALID_AUTHORIZATION
	}
	authorization = "KakaoAK" + " " + requestor.authorization

	var contentType string
	if len(requestor.contentType) == 0 {
		return ERROR_INVALID_CONTENT_TYPE
	}
	contentType = requestor.contentType

	httpRequest.Header = map[string][]string{
		"Authorization": {authorization},
		"Content-Type":  {contentType},
	}

	client := http.DefaultClient
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return err
	}

	for {
		var bytes []byte = make([]byte, 4096)
		n, err := httpResponse.Body.Read(bytes)
		if err != nil {
			httpResponse.Body.Close()
			pipeWriter.Close()
			break
		}
		pipeWriter.Write(bytes[:n])
	}

	return nil
}
