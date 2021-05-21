package ps

import (
	"bufio"
	"io"
	"unicode"
)

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buf: bufio.NewReader(r),
	}
}

type Reader struct {
	buf *bufio.Reader
}

func (reader *Reader) SkipWhite() {
	for {
		r, err := reader.readRune()
		if err != nil {
			return
		}
		// '%' から行末までコメントアウト
		if r == '%' {
			for {
				r, err = reader.readRune()
				if err != nil {
					return
				}
				if r == '\n' {
					break
				}
			}
			continue
		}
		if !unicode.IsSpace(r) {
			reader.unreadRune()
			return
		}
	}
}

func (reader *Reader) readRune() (rune, error) {
	r, _, err := reader.buf.ReadRune()
	return r, err
}

func (reader *Reader) unreadRune() error {
	err := reader.buf.UnreadRune()
	return err
}
