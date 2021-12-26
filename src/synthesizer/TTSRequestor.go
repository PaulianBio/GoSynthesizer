package synthesizer

import "io"

type TTSRequestor interface {
	RequestSynthesizer(message string) (io.Reader, error)
}
