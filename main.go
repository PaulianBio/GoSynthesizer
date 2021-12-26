package main

import (
	"bytes"
	"io"
	"log"

	synthesizer "github.com/PaulianBio/src/synthesizer/server/kakao-i"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

const KAKAOI_AUTH_KEY string = "9b8c28784daa4d246244263249e7f956"

func main() {
	pr, pw := io.Pipe()

	go func() {
		messages := []synthesizer.KakaoVoiceMessage{
			{synthesizer.KAKAOI_VOICE_STYLE_WOMAN_READ_CALM, "여성 차분한 낭독체로 발화합니다."},
			{synthesizer.KAKAOI_VOICE_STYLE_WOMAN_DIALOG_BRIGHT, "여성 밝은 대화체로 발화합니다."},
			{synthesizer.KAKAOI_VOICE_STYLE_MAN_READ_CALM, "남성 차분한 낭독체로 발화합니다."},
			{synthesizer.KAKAOI_VOICE_STYLE_MAN_DIALOG_BRIGHT, "남성 밝은 대화체로 발화합니다."},
		}

		requestor := synthesizer.CreateKakaoTTSRequestor(KAKAOI_AUTH_KEY)
		err := requestor.RequestSynthesizerAsync(messages, pw)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	done := make(chan error)
	go func() {
		byteBuf := &bytes.Buffer{}
		for {
			var buffer []byte = make([]byte, 4096)
			n, err := pr.Read(buffer)
			if err != nil {
				break
			}
			byteBuf.Write(buffer[:n])
		}
		log.Println("Voice message synthesized")

		dec, data, err := minimp3.DecodeFull(byteBuf.Bytes())
		defer dec.Close()
		if err != nil {
			done <- err
			return
		}

		ctx, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 4096)
		defer ctx.Close()
		if err != nil {
			done <- err
			return
		}

		player := ctx.NewPlayer()
		defer player.Close()
		if err != nil {
			done <- err
			return
		}
		player.Write(data)
		log.Println("Voice message played", len(data))

		done <- nil
	}()

	err := <-done
	if err != nil {
		log.Fatalln(err)
	}
}
