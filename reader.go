package hexunescape

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

func invalidEscapeError(ch byte, r *bufio.Reader, err error) error {
	in := make([]byte, 16)
	r.Read(in)
	return fmt.Errorf("invalid escape sequence: %c%s; %v", ch, in, err)
}

func parseHexEscape(r *bufio.Reader) (byte, error) {
	ch, err := r.ReadByte()
	if ch != 'x' || err != nil {
		return 0, invalidEscapeError(ch, r, err)
	}
	encoded := make([]byte, 2)
	for i := 0; i < 2; i++ {
		if encoded[i], err = r.ReadByte(); err != nil {
			return 0, invalidEscapeError(ch, r, err)
		}
	}
	decoded := make([]byte, 1)
	if n, err := hex.Decode(decoded, encoded); err != nil || n != 1 {
		return 0, invalidEscapeError(ch, r, err)
	}
	return decoded[0], nil
}

type reader struct {
	r *bufio.Reader
}

func (r *reader) Read(p []byte) (int, error) {
	max := len(p) - 1
	i := 0
	for {
		ch, err := r.r.ReadByte()
		if err != nil {
			return i, err
		}
		if ch == '\\' {
			if ch, err = parseHexEscape(r.r); err != nil {
				return i, err
			}
		}
		p[i] = ch
		i++
		if i == max || ch == '\n' {
			break
		}
	}
	return i, nil
}

// NewReader constructs a new hexadecimal escaped stream reader.
func NewReader(r io.Reader) io.Reader {
	return &reader{r: bufio.NewReader(r)}
}
