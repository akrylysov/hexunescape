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

// Unescape decodes hexadecimal escaped input text from io.Reader into io.Writer.
func Unescape(dst io.Writer, src io.Reader) error {
	bdst := bufio.NewWriter(dst)
	bsrc := bufio.NewReader(src)
	for {
		ch, err := bsrc.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if ch == '\\' {
			if ch, err = parseHexEscape(bsrc); err != nil {
				return err
			}
		}
		if err := bdst.WriteByte(ch); err != nil {
			return err
		}
		if ch == '\n' {
			if err := bdst.Flush(); err != nil {
				return err
			}
		}
	}
	return bdst.Flush()
}
