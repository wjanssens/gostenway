package rtxt

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type ReliableTxtEncoding int

const (
	Utf8         ReliableTxtEncoding = 0
	Utf16                            = 1
	Utf16Reverse                     = 2
	Utf32                            = 3
)

func Join(lines []string) string {
	return strings.Join(lines, "\n")
}
func Split(s string) []string {
	return strings.Split(s, "\n")
}

func Encoder(w io.Writer, enc ReliableTxtEncoding) (io.Writer, error) {
	switch enc {
	case Utf32:
		return nil, fmt.Errorf("UTF32 encoding not implemented")
	case Utf16:
		t := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		w = transform.NewWriter(w, t.NewEncoder().Transformer)
	case Utf16Reverse:
		t := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		w = transform.NewWriter(w, t.NewEncoder().Transformer)
	default:
		t := unicode.UTF8BOM
		w = transform.NewWriter(w, t.NewEncoder().Transformer)
	}
	return w, nil
}

func WriteLines(w io.Writer, lines []string, enc ReliableTxtEncoding) (n int, err error) {
	var r int
	e, err := Encoder(w, enc)
	if err != nil {
		return 0, err
	}
	for l := range lines {
		if l > 0 {
			if n, err := io.WriteString(e, "\n"); err == nil {
				return n + r, err
			} else {
				n += r
			}
		}
		if n, err = io.WriteString(e, lines[l]); err == nil {
			return n + r, err
		} else {
			n += r
		}
	}
	return r, nil
}

func ReadLines(r io.Reader) (lines []string, err error) {
	res := make([]string, 0)
	s := ScanLines(r)
	for s.Scan() {
		res = append(res, s.Text())
	}
	return res, s.Err()
}

func ScanLines(r io.Reader) *bufio.Scanner {
	t := unicode.BOMOverride(unicode.UTF8BOM.NewDecoder())
	tr := transform.NewReader(r, t)
	s := bufio.NewScanner(tr)
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, data, bufio.ErrFinalToken
		}
		if i := strings.Index(string(data), "\n"); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, bufio.ErrFinalToken
		}
		return
	})
	return s
}
