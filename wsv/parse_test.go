package wsv

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	type test struct {
		input    string
		preserve string
		minimal  string
	}
	tests := []test{
		// {"", "", ""},
		// {" ", " ", ""},
		// {"  ", "  ", ""},
		// {"a", "a", "a"},
		// {"a ", "a ", "a"},
		// {"a  ", "a  ", "a"},
		// {" a", " a", "a"},
		// {"  a", "  a", "a"},
		// {"  a  ", "  a  ", "a"},
		// {"a b", "a b", "a b"},
		// {"a  b", "a  b", "a b"},
		// {" a b", " a b", "a b"},
		// {"  a b", "  a b", "a b"},
		// {"  a  b", "  a  b", "a b"},
		// {"a b ", "a b ", "a b"},
		// {"a  b  ", "a  b  ", "a b"},
		// {" a b ", " a b ", "a b"},
		// {"  a b ", "  a b ", "a b"},
		// {"  a  b  ", "  a  b  ", "a b"},
		{"#", "#", ""},
		// {" #", " #", ""},
		// {"  #", "  #", ""},
		// {"a#", "a#", "a"},
		// {"a #", "a #", "a"},
		// {"a  #", "a  #", "a"},
		// {" a#", " a#", "a"},
		// {"  a#", "  a#", "a"},
		// {"  a  #", "  a  #", "a"},
		// {"a b#", "a b#", "a b"},
		// {"a  b#", "a  b#", "a b"},
		// {" a b#", " a b#", "a b"},
		// {"  a b#", "  a b#", "a b"},
		// {"  a  b#", "  a  b#", "a b"},
		// {"a b #", "a b #", "a b"},
		// {"a  b  #", "a  b  #", "a b"},
		// {" a b #", " a b #", "a b"},
		// {"  a b #", "  a b #", "a b"},
		// {"  a  b  #", "  a  b  #", "a b"},
		// {"#c", "#c", ""},
		// {" #c", " #c", ""},
		// {"  #c", "  #c", ""},
		// {"a#c", "a#c", "a"},
		// {"a #c", "a #c", "a"},
		// {"a  #c", "a  #c", "a"},
		// {" a#c", " a#c", "a"},
		// {"  a#c", "  a#c", "a"},
		// {"  a  #c", "  a  #c", "a"},
		// {"a b#c", "a b#c", "a b"},
		// {"a  b#c", "a  b#c", "a b"},
		// {" a b#c", " a b#c", "a b"},
		// {"  a b#c", "  a b#c", "a b"},
		// {"  a  b#c", "  a  b#c", "a b"},
		// {"a b #c", "a b #c", "a b"},
		// {"a  b  #c", "a  b  #c", "a b"},
		// {" a b #c", " a b #c", "a b"},
		// {"  a b #c", "  a b #c", "a b"},
		// {"  a  b  #c", "  a  b  #c", "a b"},
		// {"\uD834\uDD1E", "\uD834\uDD1E", "\uD834\uDD1E"},
		// {"#\uD834\uDD1E", "#\uD834\uDD1E", ""},
		// {`""`, `""`, `""`},
		// {`"" `, `"" `, `""`},
		// {`"\uD834\uDD1E"`, `\uD834\uDD1E`, `\uD834\uDD1E`},
		// {"-", "-", "-"},
		// {"-a", "-a", "-a"},
	}

	for i, test := range tests {
		if l, err := ParseLine(test.input, true); err == nil {
			fmt.Println(l.comment)

			preserve := l.String()
			if preserve != test.preserve {
				t.Errorf("%v: expected %v, got %v", i, test.preserve, preserve)
			}

			minimal := l.ValuesString()
			if minimal != test.minimal {
				t.Errorf("%v: expected %v, got %v", i, test.minimal, minimal)
			}

		} else {
			t.Errorf("%v: unexpected error %v", i, err)
		}
		if l, err := ParseLine(test.input, false); err == nil {
			preserve := l.String()
			if preserve != test.minimal {
				t.Errorf("%v: expected %v, got %v", i, test.minimal, preserve)
			}

			minimal := l.ValuesString()
			if minimal != test.minimal {
				t.Errorf("%v: expected %v, got %v", i, test.minimal, minimal)
			}
		} else {
			t.Errorf("%v: unexpected error %v", i, err)
		}
	}
}
