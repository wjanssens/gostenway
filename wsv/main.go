package wsv

import (
	"bytes"
	"math/big"
	"strings"
)

func SerializeValue(s string, isNull bool) string {
	if isNull {
		return "-"
	} else if len(s) == 0 {
		return "\"\""
	} else if s == "-" {
		return "\"-\""
	} else if ContainsSpecialChar(s) {
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
	hash    bool
	comment string
}

func NewLine() Line {
	return Line{}
}
func (l *Line) Len() int {
	return len(l.values)
}
func (l *Line) HasValues() bool {
	return l.values != nil && len(l.values) > 0
}
func (l *Line) GetValues() []string {
	return l.values
}
func (l *Line) SetValues(values []string) {
	l.values = values
}
func (l *Line) SetValue(i int, value string) {
	l.values[i] = value
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
func (l *Line) HasSpaces() bool {
	return l.spaces != nil && len(l.spaces) > 0
}
func (l *Line) GetSpaces() []string {
	return l.spaces
}
func (l *Line) SetSpaces(spaces []string) error {
	if err := ValidateSpaces(spaces); err != nil {
		return err
	}
	l.spaces = spaces
	return nil
}
func (l *Line) SetSpace(i int, space string) error {
	if err := ValidateSpace(space, i == 0); err != nil {
		return err
	}
	l.spaces[i] = space
	return nil
}
func (l *Line) HasComment() bool {
	return l.hash && l.comment != ""
}
func (l *Line) GetComment() (string, bool) {
	return l.comment, l.hash
}
func (l *Line) SetComment(s string) error {
	if err := ValidateComment(s); err != nil {
		return err
	}
	l.hash = true
	l.comment = s
	return nil
}
func (l *Line) ClearComment() {
	l.hash = false
	l.comment = ""
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
		result = append(result, SerializeValue(v, null))
	}
	if spacect > valuect {
		result = append(result, l.spaces[valuect])
	}
	if l.hash {
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
		result = append(result, SerializeValue(v, null))
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
