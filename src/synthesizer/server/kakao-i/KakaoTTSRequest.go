package synthesizer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	synthesizer "github.com/PaulianBio/src/synthesizer/server"
)

// <-- KAKAOI REQUEST API -->
const KAKAOI_METHOD string = http.MethodPost
const KAKAOI_SCHEME string = "https"
const KAKAOI_HOST string = "kakaoi-newtone-openapi.kakao.com"
const KAKAOI_PATH string = "/v1/synthesize"
const KAKAOI_CONTENT_TYPE string = "application/xml"

// <-- VOICE STYLE -->
const KAKAOI_VOICE_STYLE_WOMAN_READ_CALM = "WOMAN_READ_CALM"
const KAKAOI_VOICE_STYLE_MAN_READ_CALM = "MAN_READ_CALM"
const KAKAOI_VOICE_STYLE_WOMAN_DIALOG_BRIGHT = "WOMAN_DIALOG_BRIGHT"
const KAKAOI_VOICE_STYLE_MAN_DIALOG_BRIGHT = "MAN_DIALOG_BRIGHT"

type KakaoTTSRequest struct {
	synthesizer.RequestApi
}

type VoiceMessage struct {
	XMLName xml.Name `xml:"voice"`
	Message string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
}

type SpeakMessage struct {
	XMLName  xml.Name `xml:"speak"`
	Messages []VoiceMessage
}

func (m *SpeakMessage) AddDefaultMessage(message string) {
	m.Messages = append(m.Messages, VoiceMessage{
		Name:    KAKAOI_VOICE_STYLE_WOMAN_READ_CALM, // default
		Message: message,
	})
}

func (m *SpeakMessage) AddCustomMessage(message string, style string) {
	m.Messages = append(m.Messages, VoiceMessage{
		Name:    style,
		Message: message,
	})
}

func (m *SpeakMessage) ParseXMLRequest() []byte {
	if len(m.Messages) == 0 {
		return nil
	}
	var marshal []byte
	marshal, _ = xml.Marshal(m)
	return marshal
}

func (request *KakaoTTSRequest) NewHttpRequest(messages *[]KakaoVoiceMessage) (*http.Request, error) {
	httpRequest, err := request.RequestApi.NewHttpRequest()
	if err != nil {
		return httpRequest, err
	}

	var speakMessage SpeakMessage
	for _, message := range *messages {
		speakMessage.AddCustomMessage(message.Message, message.Style)
	}

	xmlRequest := bytes.NewReader(speakMessage.ParseXMLRequest())
	httpRequest, _ = http.NewRequest(KAKAOI_METHOD, fmt.Sprintf("%s://%s%s", KAKAOI_SCHEME, KAKAOI_HOST, KAKAOI_PATH), xmlRequest)
	return httpRequest, nil
}
