package wsv

type LineBuilder struct {
	Line   *Line
	Errors []error
}

func NewLineBuilder() *LineBuilder {
	return &LineBuilder{
		Line: NewLine(),
	}
}
func (b *LineBuilder) Value(i int, value string) *LineBuilder {
	b.Line.SetValue(i, value)
	return b
}
func (b *LineBuilder) Values(values []string) *LineBuilder {
	b.Line.SetValues(values)
	return b
}
func (b *LineBuilder) Space(i int, space string) *LineBuilder {
	if err := b.Line.SetSpace(i, space); err != nil {
		b.Errors = append(b.Errors, err)
	}
	return b
}
func (b *LineBuilder) Spaces(spaces []string) *LineBuilder {
	if err := b.Line.SetSpaces(spaces); err != nil {
		b.Errors = append(b.Errors, err)
	}
	return b
}
func (b *LineBuilder) Comment(comment string) *LineBuilder {
	if err := b.Line.SetComment(comment); err != nil {
		b.Errors = append(b.Errors, err)
	}
	return b
}
func (b *LineBuilder) ClearComment() *LineBuilder {
	b.Line.ClearComment()
	return b
}
func (b *LineBuilder) Nil(i int) *LineBuilder {
	b.Line.SetNil(i)
	return b
}
func (b *LineBuilder) NotNil(i int) *LineBuilder {
	b.Line.UnsetNil(i)
	return b
}
func (b *LineBuilder) Build() (*Line, []error) {
	return b.Line, b.Errors
}
