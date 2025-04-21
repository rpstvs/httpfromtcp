package response

import (
	"fmt"
	"io"

	"github.com/rpstvs/httpfromtcp/internal/headers"
)

type writeState int

const (
	writerStateStatusLine writeState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writerState writeState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}

	defer func() {
		w.writerState = writerStateHeaders
	}()

	_, err := w.writer.Write(getStatusLine(code))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.writerState)
	}

	defer func() { w.writerState = writerStateBody }()

	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	chunksize := len(p)

	nTotal := 0

	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunksize)

	if err != nil {
		return nTotal, err
	}

	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))

	if err != nil {
		return nTotal, err
	}

	nTotal += n

	return nTotal, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cant write body on state %d", w.writerState)
	}
	n, err := w.writer.Write([]byte("0\r\n\r\n"))

	if err != nil {
		return n, err
	}

	return n, nil
}
