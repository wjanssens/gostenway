package wsv

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	type test struct {
		builder  *LineBuilder
		preserve string
		minimal  string
	}
	v1 := []string{"a"}
	v2 := []string{"a", "b"}
	s1 := []string{"\t\t"}
	s1e := []string{""}
	s2 := []string{"\t\t", "  "}
	tests := []test{
		{NewLineBuilder(), "", ""},
		{NewLineBuilder().Values(v1), "a", "a"},
		{NewLineBuilder().Values(v2), "a b", "a b"},
		{NewLineBuilder().Spaces(s1), "\t\t", ""},
		{NewLineBuilder().Spaces(s2), "\t\t", ""},
		{NewLineBuilder().Values(v1).Spaces(s1), "\t\ta", "a"},
		{NewLineBuilder().Values(v1).Spaces(s2), "\t\ta  ", "a"},
		{NewLineBuilder().Comment("c"), "#c", ""},
		{NewLineBuilder().Comment("c").Spaces(s1e), "#c", ""},
		{NewLineBuilder().Comment("c").Spaces(s1), "\t\t#c", ""},
		{NewLineBuilder().Comment("c").Spaces(s2), "\t\t#c", ""},
		{NewLineBuilder().Values(v1).Comment("c"), "a#c", "a"},
		{NewLineBuilder().Values(v1).Spaces(s1e), "a", "a"},
		{NewLineBuilder().Values(v1).Spaces(s1).Comment("c"), "\t\ta#c", "a"},
		{NewLineBuilder().Values(v1).Spaces(s2).Comment("c"), "\t\ta  #c", "a"},
		{NewLineBuilder().Nil(0), "-", "-"},
		{NewLineBuilder().Values(v1).Nil(1), "a -", "a -"},
	}

	for i, test := range tests {
		if line, err := test.builder.Build(); err != nil {
			t.Errorf("%v: unexpected errors: %v", i, err)
		} else {
			preserve := line.String()
			if preserve != test.preserve {
				t.Errorf("%v: expected %v, got %v", i, test.preserve, preserve)
			}
			minimal := line.ValuesString()
			if minimal != test.minimal {
				t.Errorf("%v: expected %v, got %v", i, test.minimal, minimal)
			}
		}
	}
}
