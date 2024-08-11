package rtxt

import (
	"bytes"
	"encoding/hex"
	"io"
	"strings"
	"testing"
)

type Pair[T any] struct {
	val T
	hex string
}

func TestEncodeUtf8(t *testing.T) {
	for _, p := range []Pair[string]{
		{"", ""},
		{"a", "61"},
		{"a¥", "61C2A5"},
		{"\uFEFF", "EFBBBF"},
		{"\uFEFF\uFEFF", "EFBBBFEFBBBF"},
		{"a\u6771", "61E69DB1"},
		{"\u0000", "00"},
	} {
		var buf bytes.Buffer
		bom, _ := hex.DecodeString("efbbbf")
		str, _ := hex.DecodeString(p.hex)
		expected := append(bom[:], str[:]...)

		w, _ := Encoder(&buf, Utf8)
		_, err := io.WriteString(w, p.val)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !bytes.Equal(expected, buf.Bytes()) {
			t.Errorf("Incorrect string; expected %v, received %v", expected, buf.Bytes())
		}
	}
}

func TestEncodeUtf16(t *testing.T) {
	for _, p := range []Pair[string]{
		{"", ""},
		{"a", "0061"},
		{"a¥", "006100A5"},
		{"\uFEFF", "FEFF"},
		{"\uFEFF\uFEFF", "FEFFFEFF"},
		{"a\u6771", "00616771"},
		{"\u0000", "0000"},
	} {
		var buf bytes.Buffer
		bom, _ := hex.DecodeString("feff")
		str, _ := hex.DecodeString(p.hex)
		expected := append(bom[:], str[:]...)

		w, _ := Encoder(&buf, Utf16)
		_, err := io.WriteString(w, p.val)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !bytes.Equal(expected, buf.Bytes()) {
			t.Errorf("Incorrect string; expected %v, received %v", expected, buf.Bytes())
		}
	}
}

func TestEncodeUtf16Reverse(t *testing.T) {
	for _, p := range []Pair[string]{
		{"", ""},
		{"a", "6100"},
		{"a¥", "6100A500"},
		{"\uFEFF", "FFFE"},
		{"\uFEFF\uFEFF", "FFFEFFFE"},
		{"a\u6771", "61007167"},
		{"\u0000", "0000"},
	} {
		var buf bytes.Buffer
		bom, _ := hex.DecodeString("fffe")
		str, _ := hex.DecodeString(p.hex)
		expected := append(bom[:], str[:]...)

		w, _ := Encoder(&buf, Utf16Reverse)
		_, err := io.WriteString(w, p.val)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !bytes.Equal(expected, buf.Bytes()) {
			t.Errorf("Incorrect string; expected %v, received %v", expected, buf.Bytes())
		}
	}
}

func TestScanUtf8(t *testing.T) {
	for _, p := range []Pair[string]{
		// {"", ""},
		{"a", "61"},
		{"a¥", "61C2A5"},
		{"\uFEFF", "EFBBBF"},
		{"\uFEFF\uFEFF", "EFBBBFEFBBBF"},
		{"a\u6771", "61E69DB1"},
		{"\u0000", "00"},
	} {
		bom, _ := hex.DecodeString("efbbbf")
		str, _ := hex.DecodeString(p.hex)
		buf := bytes.NewBuffer(append(bom[:], str[:]...))

		arr, err := ReadLines(buf)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.EqualFold(arr[0], p.val) {
			t.Errorf("Incorrect string; expected %v, received %v", p.val, arr[0])
		}
	}
}

func TestScanUtf16(t *testing.T) {
	for _, p := range []Pair[string]{
		// {"", ""},
		{"a", "0061"},
		{"a¥", "006100A5"},
		{"\uFEFF", "FEFF"},
		{"\uFEFF\uFEFF", "FEFFFEFF"},
		{"a\u6771", "00616771"},
		{"\u0000", "0000"},
	} {
		bom, _ := hex.DecodeString("feff")
		str, _ := hex.DecodeString(p.hex)
		buf := bytes.NewBuffer(append(bom[:], str[:]...))

		arr, err := ReadLines(buf)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.EqualFold(arr[0], p.val) {
			t.Errorf("Incorrect string; expected %v, received %v", p.val, arr[0])
		}
	}
}

func TestScanUtf16Reverse(t *testing.T) {
	for _, p := range []Pair[string]{
		// {"", ""},
		{"a", "6100"},
		{"a¥", "6100A500"},
		{"\uFEFF", "FFFE"},
		{"\uFEFF\uFEFF", "FFFEFFFE"},
		{"a\u6771", "61007167"},
		{"\u0000", "0000"},
	} {
		bom, _ := hex.DecodeString("fffe")
		str, _ := hex.DecodeString("0a")
		buf := bytes.NewBuffer(append(bom[:], str[:]...))

		arr, err := ReadLines(buf)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.EqualFold(arr[0], p.val) {
			t.Errorf("Incorrect string; expected %v, received %v", p.val, arr[0])
		}
	}
}

func TestScanLinesUtf8(t *testing.T) {
	for _, p := range []Pair[[]string]{
		{[]string{}, ""},
		{[]string{"", ""}, "0A"},
		{[]string{"A", "B"}, "410a42"},
		{[]string{"\u0000", "\u0000"}, "000A00"},
		{[]string{"𝄞", "𝄞"}, "f09d849e0Af09d849e"},
	} {
		bom, _ := hex.DecodeString("efbbbf")
		str, _ := hex.DecodeString(p.hex)
		buf := bytes.NewBuffer(append(bom[:], str[:]...))

		arr, err := ReadLines(buf)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		for i, str := range p.val {
			if !strings.EqualFold(arr[i], str) {
				t.Errorf("Incorrect string; expected %v, received %v", str, arr[i])
			}
		}
	}
}
