package wsv

import (
	"fmt"
	"slices"
)

var ws = []rune{
	0x0009,
	0x000B,
	0x000C,
	0x000D,
	0x0020,
	0x0085,
	0x00A0,
	0x1680,
	0x2000,
	0x2001,
	0x2002,
	0x2003,
	0x2004,
	0x2005,
	0x2006,
	0x2007,
	0x2008,
	0x2009,
	0x200A,
	0x2028,
	0x2029,
	0x202F,
	0x205F,
	0x3000,
}

func isWs(r rune) bool {
	return slices.Contains(ws, r)
}

func ValidateSpaces(spaces []string) error {
	for i, s := range spaces {
		if err := ValidateSpace(s, i == 0); err != nil {
			return err
		}
	}
	return nil
}
func ValidateSpace(s string, first bool) error {
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

func ValidateComment(s string) error {
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

func ContainsSpecialChar(s string) bool {
	for _, r := range s {
		if r == 0x0022 || r == 0x0023 || r == 0x000a || isWs(r) {
			return true
		}
	}
	return false
}

func IsSpecial(s string) bool {
	return s == "" || s == "-" || ContainsSpecialChar(s)
}
