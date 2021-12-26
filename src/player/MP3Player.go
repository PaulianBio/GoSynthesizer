package player

import (
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

type MP3Player struct{}

func (p *MP3Player) Play(reader io.Reader) error {
	log.Println("MP3Player::Play()")

	dec, err := minimp3.NewDecoder(reader)
	if err != nil {
		return err
	}
	defer dec.Close()
	<-dec.Started()

	ctx, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 4096)
	if err != nil {
		return err
	}
	defer ctx.Close()

	done := make(chan bool)
	go func() {
		player := ctx.NewPlayer()
		for {
			var data []byte = make([]byte, 1024)
			n, err := dec.Read(data)
			if err != nil {
				break
			}
			player.Write(data[:n])
			log.Println("Player writes bytes", n)
		}
		player.Close()
		player = nil
		done <- true
	}()

	<-done
	time.Sleep(time.Duration(500) * time.Millisecond)
	return nil
}

func (p *MP3Player) PlayAll(bytes []byte) error {
	log.Println("MP3Player::PlayAll()")

	dec, data, err := minimp3.DecodeFull(bytes)
	if err != nil {
		return err
	}
	defer dec.Close()

	ctx, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 4096)
	if err != nil {
		return err
	}
	defer ctx.Close()

	player := ctx.NewPlayer()
	if err != nil {
		return err
	}
	defer player.Close()
	player.Write(data)

	log.Println("Voice message played:", len(data))
	return nil
}
