package hexunescape

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
)

func TestUnescape(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{`\x22`, `"`},
		{`\x22\x22\x22`, `"""`},
		{`a\x22`, `a"`},
		{`\x22a`, `"a`},
		{`a\x22b`, `a"b`},
		{`a\x22b\x22`, `a"b"`},
		{
			"\xd0\x94\xd1\x83\xd0\xb1\xd1\x8a \xe2\x80\x94 \xd0\xb4\xd0\xb5\xd1\x80\xd0\xb5\xd0\xb2\xd0\xbe",
			"Дубъ — дерево",
		},
	}
	for _, testCase := range testCases {
		w := bytes.NewBuffer(nil)
		if err := Unescape(w, strings.NewReader(testCase.in)); err != nil {
			t.Errorf("%s err: %v", testCase.in, err)
		}
		if out := w.String(); out != testCase.out {
			t.Errorf("%s got: %s want: %s", testCase.in, out, testCase.out)
		}
	}
}

func TestUnescapeError(t *testing.T) {
	testCases := []string{
		`\`,
		`\\x22`,
		`\x2`,
		`\xzza`,
	}
	for _, testCase := range testCases {
		if err := Unescape(ioutil.Discard, strings.NewReader(testCase)); err == nil {
			t.Errorf("%s got: %v want: error", testCase, err)
		}
	}
}

func generateRandomData(n int) ([]byte, []byte) {
	unescaped := bytes.NewBuffer(make([]byte, n))
	escaped := bytes.NewBuffer(make([]byte, n))
	for i := 0; i < n; i++ {
		b := byte(rand.Intn(255))
		unescaped.WriteByte(b)
		if b < 'A' || b > 'z' || b == '\\' {
			escaped.Write([]byte{'\\', 'x'})
			be := make([]byte, 2)
			hex.Encode(be, []byte{b})
			escaped.Write(be)
		} else {
			escaped.WriteByte(b)
		}
	}
	return unescaped.Bytes(), escaped.Bytes()
}

func TestUnescapeRandom(t *testing.T) {
	unescaped, escaped := generateRandomData(1024 * 1024)
	w := bytes.NewBuffer(nil)
	if err := Unescape(w, bytes.NewReader(escaped)); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(unescaped, w.Bytes()) {
		t.Error()
	}
}

func benchmarkUnescape(b *testing.B, inputLen int) {
	_, escaped := generateRandomData(inputLen)
	b.ResetTimer()
	b.SetBytes(int64(len(escaped)))
	for i := 0; i < b.N; i++ {
		if err := Unescape(ioutil.Discard, bytes.NewReader(escaped)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkUnescape1KB(b *testing.B) { benchmarkUnescape(b, 1024) }
func BenchmarkUnescape8KB(b *testing.B) { benchmarkUnescape(b, 8*1024) }
func BenchmarkUnescape1MB(b *testing.B) { benchmarkUnescape(b, 1024*1024) }
func BenchmarkUnescape8MB(b *testing.B) { benchmarkUnescape(b, 8*1024*1024) }
