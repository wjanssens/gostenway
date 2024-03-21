package sml

import (
	"io"
	"strings"

	"github.com/wjanssens/rtxt"
	"github.com/wjanssens/wsv"
)

func Parse(r io.Reader, preserveWhitespaceAndComments bool, lineIndexOffset int) (*Node, error) {
	root := NewRoot()
	curr := &root
	stk := stack{}
	stk.Push(curr)

	s := rtxt.ScanLines(r)

	lineIndex := lineIndexOffset - 1
	for s.Scan() {
		lineIndex++
		if lineIndex < lineIndexOffset {
			continue
		}
		if line, err := wsv.ParseLine(s.Text(), preserveWhitespaceAndComments); err != nil {
			return curr, err // TODO better to return root, or current node?
		} else {
			values := line.GetValues()
			name := values[0]
			if len(values) == 0 {
				n := Node{start: &line}
				curr.children = append(curr.children, n)
			} else if len(values) == 1 {
				if strings.EqualFold(name, "end") {
					curr := stk.Pop()
					curr.end = &line
				} else {
					n := Node{start: &line, children: make([]Node, 0)}
					curr.children = append(curr.children, n)
				}
			} else {
				n := Node{}
				curr.children = append(curr.children, n)
			}
		}
	}
	return &root, nil
}
