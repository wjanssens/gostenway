package wsv

import (
	"testing"
)

const testTable = `a 	U+0061    61            0061        "Latin Small Letter A"
~ 	U+007E    7E            007E        Tilde
¥ 	U+00A5    C2_A5         00A5        "Yen Sign"
» 	U+00BB    C2_BB         00BB        "Right-Pointing Double Angle Quotation Mark"
½ 	U+00BD    C2_BD         00BD        "Vulgar Fraction One Half"
¿ 	U+00BF    C2_BF         00BF        "Inverted Question Mark"
ß 	U+00DF    C3_9F         00DF        "Latin Small Letter Sharp S"
ä 	U+00E4    C3_A4         00E4        "Latin Small Letter A with Diaeresis"
ï 	U+00EF    C3_AF         00EF        "Latin Small Letter I with Diaeresis"
œ 	U+0153    C5_93         0153        "Latin Small Ligature Oe"
€ 	U+20AC    E2_82_AC      20AC        "Euro Sign"
東 	U+6771    E6_9D_B1      6771        "CJK Unified Ideograph-6771"
𝄞 	U+1D11E   F0_9D_84_9E   D834_DD1E   "Musical Symbol G Clef"
𠀇 	U+20007   F0_A0_80_87   D840_DC07   "CJK Unified Ideograph-20007"`

func TestValidateWhitespace(t *testing.T) {
	firstValid := []string{
		"",
		" ",
		"\t",
		"  \t  ",
		"\u0009\u000B\u000C\u000D\u0020\u0085\u00A0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A\u2028\u2029\u202F\u205F\u3000",
	}
	nonfirstValid := []string{
		" ",
		"\t",
		"  \t  ",
		"\u0009\u000B\u000C\u000D\u0020\u0085\u00A0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A\u2028\u2029\u202F\u205F\u3000",
	}
	firstInvalid := []string{
		"a",
	}
	nonfirstInvalid := []string{
		"",
		"a",
		"  a",
	}
	for i, s := range firstValid {
		if err := ValidateSpace(s, true); err != nil {
			t.Errorf("%v: expected %v to be valid: %v", i, s, err)
		}
	}
	for i, s := range firstInvalid {
		if err := ValidateSpace(s, true); err == nil {
			t.Errorf("%v: expected %v to be invalid", i, s)
		}
	}
	for i, s := range nonfirstValid {
		if err := ValidateSpace(s, false); err != nil {
			t.Errorf("%v: expected %v to be valid: %v", i, s, err)
		}
	}
	for i, s := range nonfirstInvalid {
		if err := ValidateSpace(s, false); err == nil {
			t.Errorf("%v: expected %v to be invalid", i, s)
		}
	}

}

func TestValidateSpaces(t *testing.T) {
	valid := [][]string{
		{},
		{" "},
		{" ", " "},
		{""},
		{"", " "},
		{"", " ", "   \t  "},
	}
	invalid := [][]string{
		{"", ""},
		{" ", "  a"},
	}
	for i, s := range valid {
		if err := ValidateSpaces(s); err != nil {
			t.Errorf("%v: expected %v to be valid: %v", i, s, err)
		}
	}
	for i, s := range invalid {
		if err := ValidateSpaces(s); err == nil {
			t.Errorf("%v: expected %v to be invalid", i, s)
		}
	}
}

func TestValidateComment(t *testing.T) {
	valid := []string{
		"",
		" ",
		"a",
		"comment",
		"#",
		"######",
		// "\uD834\uDD1E",
	}
	invalid := []string{
		"\n",
		// "\uD834",
		// "\uD834\uD834",
		// "\uDD1E",
		// "\uDD1E\uDD1E",
	}
	for i, s := range valid {
		if err := ValidateComment(s); err != nil {
			t.Errorf("%v: expected %v to be valid: %v", i, s, err)
		}
	}
	for i, s := range invalid {
		if err := ValidateComment(s); err == nil {
			t.Errorf("%v: expected %v to be invalid", i, s)
		}
	}
}

func TestHasValues(t *testing.T) {
	type test struct {
		values   []string
		hasValue bool
	}
	tests := []test{
		{nil, false},
		{[]string{""}, true},
		{[]string{"v"}, true},
	}
	for i, p := range tests {
		if l, err := NewLineBuilder().Values(p.values).Build(); err == nil {
			if p.hasValue != l.HasValues() {
				t.Errorf("%v: expected %v", i, p.hasValue)
			}
		} else {
			t.Errorf("%v: unexpected error %v", i, err)
		}
	}
}

func TestHasSpaces(t *testing.T) {
	type test struct {
		spaces    []string
		hasSpaces bool
	}
	tests := []test{
		{[]string{}, false},
		{[]string{""}, true},
		{[]string{" "}, true},
	}
	for i, test := range tests {
		if l, err := NewLineBuilder().Spaces(test.spaces).Build(); err == nil {
			if test.hasSpaces != l.HasSpaces() {
				t.Errorf("%v: expected %v", i, test.hasSpaces)
			}
		} else {
			t.Errorf("%v: unexpected error %v", i, err)
		}
	}
}

func TestHasComments(t *testing.T) {
	type test struct {
		comment    string
		hasComment bool
	}
	tests := []test{
		{"", false},
		{"c", true},
	}
	for i, test := range tests {
		if l, err := NewLineBuilder().Comment(test.comment).Build(); err == nil {
			if test.hasComment != l.HasComment() {
				t.Errorf("%v: expected %v", i, test.hasComment)
			}
		} else {
			t.Errorf("%v: unexpected error %v", i, err)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	valid := []string{
		"",
		"a",
		"a b",
		"a b c",
		"\t\t",
		"\t\ta",
		"\t\ta  ",
		"#c",
		"\t\t#c",
		"a#c",
		"a #c",
		"\t\ta #c",
		"\t\ta  #c  ",
		"a b#c",
		"-",
		"a -",
	}

	for i, s := range valid {
		if l, err := ParseLine(s, true); err != nil {
			t.Errorf("Expected %v: %v to be valid: %v", i, s, err)
		} else {
			t.Logf("here %v", l)
			if x := l.String(); x != s {
				t.Errorf("Expected %v: expected '%v' == '%v'", i, x, s)
			}
		}
	}

}
