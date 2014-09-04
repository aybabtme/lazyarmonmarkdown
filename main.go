package main

import (
	"errors"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	input, err := ioutil.ReadAll(timeoutReader{os.Stdin, time.Second})
	if err != nil {
		log.Fatalf("oh noes! did you pipe something to stdin? %v", err)
	}
	_, err = os.Stdout.Write(blackfriday.MarkdownCommon(input))
	if err != nil {
		log.Fatalf("oh noes! maybe it's not markdown? %v", err)
	}
}

type timeoutReader struct {
	io.Reader
	dur time.Duration
}

func (t timeoutReader) Read(p []byte) (n int, err error) {
	out := make(chan struct{})

	go func() {
		n, err = t.Reader.Read(p)
		close(out)
	}()

	select {
	case <-time.After(t.dur):
		return n, errors.New("timed out waiting to read ")
	case <-out:
		return
	}
}
