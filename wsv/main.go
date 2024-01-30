package wsv

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
)

func validateSpaces(spaces []string) error {
	for i, s := range spaces {
		if err := validateSpace(s, i == 0); err != nil {
			return err
		}
	}
	return nil
}
func validateSpace(s string, first bool) error {
	if s == "" && !first {
		return fmt.Errorf("Non-first whitespace string cannot be empty")
	}

	for i, r := range s {
		if !isWs(r) {
			return fmt.Errorf("Invalid code unit '%v' in whitespace string at index %v", r, i)
		}
	}
	return nil
}

func validateComment(s string) error {
	var ctl bool = false
	for _, r := range s {
		if ctl && r >= 0xdc00 {
			return fmt.Errorf("Invalid UTF-16 String")
		} else if r == 0x000a {
			return fmt.Errorf("Line feed in comment is not allowed")
		} else if r > 0xd800 && r <= 0xdfff {
			ctl = true
		} else if ctl {
			ctl = false
		}
	}
	if ctl {
		return fmt.Errorf("Invalid UTF-16 String")
	} else {
		return nil
	}
}

func containsSpecialChar(s string) bool {
	for _, r := range s {
		if isWs(r) {
			return true
		}
		// cr = value.charCodeAt(i)
		// if (c >= 0xD800 && c <= 0xDFFF) {
		// 	i++
		// 	if (c >= 0xDC00 || i >= value.length) { throw new InvalidUtf16StringError() }
		// 	const secondCodeUnit: number = value.charCodeAt(i)
		// 	if (!(secondCodeUnit >= 0xDC00 && secondCodeUnit <= 0xDFFF)) { throw new InvalidUtf16StringError() }
		// }
	}
	return false
}

func IsSpecial(s string) bool {
	return s == "" || s == "-" || containsSpecialChar(s)
}

func serializeValue(s string, isNull bool) string {
	if isNull {
		return "-"
	} else if len(s) == 0 {
		return "\"\""
	} else if s == "-" {
		return "\"-\""
	} else if containsSpecialChar(s) {
		var b bytes.Buffer
		b.WriteRune(0x0022) // "

		for _, r := range s {
			switch r {
			case 0x000a:
				b.WriteRune(0x0022) // "
				b.WriteRune(0x002f) // /
				b.WriteRune(0x0022) // "
			case 0x0022:
				b.WriteRune(0x0022)
				b.WriteRune(0x0022)
			default:
				b.WriteRune(r)
			}
		}
		b.WriteRune(0x022)

		return b.String()
	} else {
		return s
	}
}

type Line struct {
	// nulls and values always have the same length
	// spaces always has one more element than values (space before, spaces between values, space between value and comment)
	values  []string
	nulls   *big.Int
	spaces  []string
	comment string
}

func NewLine(values []string, spaces []string, comment string) (Line, error) {
	line := Line{}
	err := line.Set(values, spaces, comment)
	return line, err
}
func (l *Line) Set(values []string, spaces []string, comment string) error {
	l.SetValues(values)
	if err := l.SetSpaces(spaces); err != nil {
		return err
	}
	if err := l.SetComment(comment); err != nil {
		return err
	}
	return nil
}

func (l *Line) Len() int {
	return len(l.values)
}

func (l *Line) GetValues() []string {
	return l.values
}
func (l *Line) SetValues(values []string) {
	l.values = values
}
func (n *Line) IsNil(i int) bool {
	return n.nulls.Bit(i) == 1
}
func (n *Line) SetNil(i int) {
	n.nulls.SetBit(n.nulls, i, 1)
}
func (n *Line) UnsetNil(i int) {
	n.nulls.SetBit(n.nulls, i, 0)
}
func (l *Line) GetSpaces() []string {
	return l.spaces
}
func (l *Line) SetSpaces(spaces []string) error {
	if err := validateSpaces(spaces); err != nil {
		return err
	}
	l.spaces = spaces
	return nil
}
func (l *Line) HasComment() bool {
	return l.comment != ""
}
func (l *Line) GetComment() string {
	return l.comment
}
func (l *Line) SetComment(s string) error {
	if err := validateComment(s); err != nil {
		return err
	}
	l.comment = s
	return nil
}
func (l *Line) String() string {
	result := make([]string, 0)
	spacect := len(l.spaces)
	valuect := len(l.values)
	for i, v := range l.values {
		sp := ""
		if spacect > i {
			sp = l.spaces[i]
		}
		null := l.nulls.Bit(i) > 0
		result = append(result, sp)
		result = append(result, serializeValue(v, null))
	}
	if spacect > valuect {
		result = append(result, l.spaces[valuect])
	}
	if len(l.comment) > 0 {
		result = append(result, "#")
		result = append(result, l.comment)
	}

	return strings.Join(result, "")
}
func (l *Line) ValuesString() string {
	result := make([]string, 0)
	for i, v := range l.values {
		if i != 0 {
			result = append(result, " ")
		}
		null := l.nulls.Bit(i) > 0
		result = append(result, serializeValue(v, null))
	}
	return strings.Join(result, "")
}

func Serialize(lines []Line) string {
	result := make([]string, 0)
	for _, l := range lines {
		result = append(result, l.String())
	}
	return strings.Join(result, "\n")
}
