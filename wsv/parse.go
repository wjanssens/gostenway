package wsv

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/wjanssens/rtxt"
)

type state int

const (
	defaultState state = 0 // receiving anything
	commentState       = 1 // receiving comment chars
	quotedState        = 2 // receiving a quoted string
	escapeState        = 3 // just saw a quote char, next char will be a quote, a slash, or a dash
	expectState        = 4 // next char must be a quote to end a quoted slash or dash
)

func ParseLine(l string, preserveWhitespaceAndComments bool) (Line, error) {
	value := strings.Builder{}
	space := strings.Builder{}
	comment := strings.Builder{}

	values := make([]string, 0)
	spaces := make([]string, 0)

	line := Line{
		values:  values,
		spaces:  spaces,
		nulls:   big.NewInt(0),
		comment: "",
		hash:    false,
	}
	if len(l) == 0 {
		return line, nil
	}

	var state = defaultState

	for _, r := range l {
		switch state {
		case defaultState:
			if r == 0x0022 { // quote
				state = quotedState
			} else if r == 0x0023 { // hash
				line.hash = true
				state = commentState
			} else if r == 0x002d { // dash
				eov(&line, &value, true)
			} else if isWs(r) { // ws
				if value.Len() > 0 {
					eov(&line, &value, false)
				}
				space.WriteRune(r)
			} else {
				if space.Len() > 0 {
					eos(&line, &space, preserveWhitespaceAndComments)
				}
				value.WriteRune(r)
			}
		case commentState:
			comment.WriteRune(r)
		case quotedState:
			if r == 0x0022 {
				state = escapeState
			} else {
				value.WriteRune(r)
			}
		case escapeState:
			if r == 0x0022 { // quote
				// "...""..." is a quote
				value.WriteRune(r)
				state = quotedState
			} else if r == 0x0023 { // hash
				// "..."#"..." is a hash
				value.WriteRune(r)
				state = expectState
			} else if r == 0x002d { // dash
				// "..."-"..." is a dash
				value.WriteRune(r)
				state = expectState
			} else if r == 0x002f { // slash
				// "..."/"..." is a LF
				value.WriteRune(0x000a)
				state = expectState
			} else {
				// last quote wasn't an excape, it was the end of the quoted string
				value.WriteRune(0x0022)
				value.WriteRune(r)
				eov(&line, &value, false)
				state = defaultState
			}
		case expectState:
			if r == 0x0022 {
				state = quotedState
			} else {
				return line, fmt.Errorf("Invalid character after escaped character")
			}
		}
	}

	if state == quotedState {
		return line, fmt.Errorf("Quoted string not closed")
	} else {
		if value.Len() > 0 {
			eov(&line, &value, false)
		} else if space.Len() > 0 {
			eos(&line, &space, preserveWhitespaceAndComments)
		}
	}
	if preserveWhitespaceAndComments {
		line.comment = comment.String()
	}
	return line, nil
}

func eov(line *Line, value *strings.Builder, null bool) {
	// fmt.Printf("eov %v, %v, %v", line, value, null)

	if len(line.values) == 0 && len(line.spaces) == 0 {
		// ensure that there is a space before the first value
		line.spaces = append(line.spaces, "")
	}
	i := len(line.values)
	line.values = append(line.values, value.String())
	if null {
		line.nulls = line.nulls.SetBit(line.nulls, i, 1)
	}
	value.Reset()
}
func eos(line *Line, space *strings.Builder, preserveWhitespaceAndComments bool) {
	// fmt.Printf("eos %v, %v, %v", line, space, preserveWhitespaceAndComments)

	if preserveWhitespaceAndComments {
		line.spaces = append(line.spaces, space.String())
	} else {
		line.spaces = append(line.spaces, " ")
	}
	space.Reset()
}

func Parse(r io.Reader, preserveWhitespaceAndComments bool, lineIndexOffset int) ([]Line, error) {
	lines := make([]Line, 0)

	s := rtxt.ScanLines(r)

	lineIndex := lineIndexOffset - 1
	for s.Scan() {
		lineIndex++
		if lineIndex < lineIndexOffset {
			continue
		}
		if line, err := ParseLine(s.Text(), preserveWhitespaceAndComments); err != nil {
			return lines, err
		} else {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
