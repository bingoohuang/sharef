package sdp

import (
	"fmt"
	"io"
	"os"
)

const (
	OfferPrompt        = "Send this offer:"
	OfferWaitingPrompt = "Please, paste the remote offer:"
)

type STDSender struct {
	writer io.Writer
}

func SenderPipe() (io.Reader, io.Writer) {
	s := &STDSender{
		writer: os.Stdout,
	}
	return os.Stdin, s
}

func (s *STDSender) Write(p []byte) (n int, err error) {
	fmt.Printf("%s\n\n", OfferPrompt)
	n, err = s.writer.Write(p)
	if err != nil {
		return
	}

	fmt.Printf("\n%s\n\n", OfferWaitingPrompt)
	return
}
