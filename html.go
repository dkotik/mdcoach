package mdcoach

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/yuin/goldmark/v2/util"
)

var _ io.Writer = (*textWriter)(nil)

// textWriter is copied from Goldmark's html renderer
// in order to render attributes
type textWriter struct {
	w util.BufWriter
}

var htmlEscapedBytes = [5]byte{0x00, '"', '&', '<', '>'}

const textWriterJumpThreshold = 32

func (w *textWriter) Write(p []byte) (int, error) {
	l := len(p)
	if l < textWriterJumpThreshold {
		return w.writeShort(p)
	}
	var positions [len(htmlEscapedBytes)]int
	found := false
	for k, c := range htmlEscapedBytes {
		positions[k] = bytes.IndexByte(p, c)
		found = found || positions[k] != -1
	}
	if !found {
		return w.w.Write(p)
	}
	written := 0
	n := 0
	for {
		i := -1
		for _, pos := range positions {
			if pos != -1 && (i == -1 || pos < i) {
				i = pos
			}
		}
		if i == -1 {
			break
		}
		if i > n {
			wr, err := w.w.Write(p[n:i])
			written += wr
			if err != nil {
				return written, err
			}
		}
		wr, err := w.w.Write(util.EscapeHTMLByte(p[i]))
		written += wr
		if err != nil {
			return written, err
		}
		n = i + 1
		for k, pos := range positions {
			if pos != -1 && pos < n {
				positions[k] = -1
				if idx := bytes.IndexByte(p[n:], htmlEscapedBytes[k]); idx != -1 {
					positions[k] = n + idx
				}
			}
		}
	}
	if n < l {
		wr, err := w.w.Write(p[n:])
		written += wr
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func (w *textWriter) writeShort(p []byte) (int, error) {
	written := 0
	n := 0
	l := len(p)
	for i := range l {
		v := util.EscapeHTMLByte(p[i])
		if v != nil {
			wr, err := w.w.Write(p[i-n : i])
			written += wr
			if err != nil {
				return written, err
			}
			n = 0
			wr, err = w.w.Write(v)
			written += wr
			if err != nil {
				return written, err
			}
			continue
		}
		n++
	}
	if n != 0 {
		wr, err := w.w.Write(p[l-n:])
		written += wr
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func (w *textWriter) WriteByte(c byte) error {
	b := [1]byte{c}
	_, err := w.Write(b[:])
	return err
}

func (w *textWriter) WriteRune(r rune) (int, error) {
	rbuf := [4]byte{}
	n := utf8.EncodeRune(rbuf[:], r)
	return w.Write(rbuf[:n])
}

func (w *textWriter) WriteString(s string) (int, error) {
	return w.Write(util.StringToReadOnlyBytes(s))
}

func (w *textWriter) Flush() error {
	return w.w.Flush()
}
