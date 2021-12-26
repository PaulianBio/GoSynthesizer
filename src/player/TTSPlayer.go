package player

import "io"

type TTSPlayer interface {
	Play(io.Reader) error
	PlayAll([]byte) error
}
