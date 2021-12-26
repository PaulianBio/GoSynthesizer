package main

import (
	"io"
	"log"

	player "github.com/PaulianBio/src/player"
	synthesizer "github.com/PaulianBio/src/synthesizer/server"
)

func main() {
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	go func() {
		requestor := synthesizer.CreateKakaoTTSRequestor()
		messages := []synthesizer.KakaoVoiceMessage{
			{Name: synthesizer.KAKAOI_VOICE_STYLE_WOMAN_READ_CALM, Message: "여성 차분한 낭독체로 발화합니다."},
			{Name: synthesizer.KAKAOI_VOICE_STYLE_WOMAN_DIALOG_BRIGHT, Message: "여성 밝은 대화체로 발화합니다."},
			{Name: synthesizer.KAKAOI_VOICE_STYLE_MAN_READ_CALM, Message: "남성 차분한 낭독체로 발화합니다."},
			{Name: synthesizer.KAKAOI_VOICE_STYLE_MAN_DIALOG_BRIGHT, Message: "남성 밝은 대화체로 발화합니다."},
		}
		err := requestor.Synthesize(&messages, pw)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	player := &player.MP3Player{}
	err := player.Play(pr)
	if err != nil {
		log.Fatalln(err)
	}
}
