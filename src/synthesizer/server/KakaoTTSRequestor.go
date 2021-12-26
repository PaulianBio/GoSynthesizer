package synthesizer

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	HEADER_KAKAOI_AUTHORIZATION string = "Authorization"
	HEADER_KAKAOI_CONTENT_TYPE  string = "Content-Type"
)

var (
	ErrInvalidContentType   error = errors.New("invalid Content-Type")
	ErrInvalidAuthorization error = errors.New("invalid Authorization")
)

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
	ContentType   string
	Authorization string
}

type KakaoVoiceMessage struct {
	XMLName xml.Name `xml:"voice"`
	Message string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
}

type SpeakMessage struct {
	XMLName  xml.Name `xml:"speak"`
	Messages []KakaoVoiceMessage
}

func (m *SpeakMessage) addDefaultMessage(message string) {
	m.Messages = append(m.Messages, KakaoVoiceMessage{
		Name:    "WOMAN_READ_CALM", // default
		Message: message,
	})
}

func (m *SpeakMessage) addCustomMessage(message string, style string) {
	m.Messages = append(m.Messages, KakaoVoiceMessage{
		Name:    style,
		Message: message,
	})
}

func (m *SpeakMessage) xmlRequest() []byte {
	if len(m.Messages) == 0 {
		return nil
	}
	var marshal []byte
	marshal, _ = xml.Marshal(m)
	return marshal
}

func (requestor *KakaoTTSRequestor) newHttpRequest(messages *[]KakaoVoiceMessage) (*http.Request, error) {
	httpRequest, err := requestor.RequestApi.GetHttpRequest()
	if err != nil {
		return httpRequest, err
	}

	var speakMessage SpeakMessage
	speakMessage.addDefaultMessage("발화를 시작합니다.")
	for _, message := range *messages {
		if len(message.Name) == 0 {
			speakMessage.addDefaultMessage(message.Message)
		} else {
			speakMessage.addCustomMessage(message.Message, message.Name)
		}
	}

	xmlRequest := bytes.NewReader(speakMessage.xmlRequest())
	httpRequest, _ = http.NewRequest(KAKAOI_METHOD, fmt.Sprintf("%s://%s%s", KAKAOI_SCHEME, KAKAOI_HOST, KAKAOI_PATH), xmlRequest)
	return httpRequest, nil
}

func (requestor *KakaoTTSRequestor) Synthesize(messages *[]KakaoVoiceMessage, writer io.Writer) error {
	var httpRequest *http.Request
	var err error

	httpRequest, err = requestor.newHttpRequest(messages)
	if err != nil {
		return err
	}

	var authorization string
	if len(requestor.Authorization) == 0 {
		return ErrInvalidAuthorization
	}
	authorization = "KakaoAK" + " " + requestor.Authorization

	var contentType string
	if len(requestor.ContentType) == 0 {
		return ErrInvalidContentType
	}
	contentType = requestor.ContentType

	httpRequest.Header = map[string][]string{
		HEADER_KAKAOI_AUTHORIZATION: {authorization},
		HEADER_KAKAOI_CONTENT_TYPE:  {contentType},
	}

	client := http.DefaultClient
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return err
	}

	log.Println("Voice message synthesized")
	for {
		var bytes []byte = make([]byte, 2048)
		n, err := httpResponse.Body.Read(bytes)
		if err != nil {
			httpResponse.Body.Close()
			break
		}
		writer.Write(bytes[:n])
	}

	return nil
}
