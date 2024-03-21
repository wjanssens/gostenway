package sml

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/wjanssens/wsv"
)

type Node struct {
	start    *wsv.Line
	end      *wsv.Line
	children []Node
}

func NewRoot() Node {
	return Node{children: make([]Node, 0)}
}

func (n *Node) IsRoot() bool {
	return n.start == nil && n.end == nil
}
func (n *Node) IsElement() bool {
	return n.start != nil && n.end != nil
}
func (n *Node) IsAttribute() bool {
	return n.start != nil && n.end == nil && n.start.HasValues()
}
func (n *Node) IsEmpty() bool {
	return n.start != nil && n.end == nil && !n.start.HasValues()
}
func (n *Node) FilterAttributes(name string) ([]Node, error) {
	if n.IsElement() || n.IsRoot() {
		result := make([]Node, 0)
		for _, child := range n.children {
			if child.IsAttribute() && strings.EqualFold(child.GetName(), name) {
				result = append(result, child)
			}
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) FilterElements(name string) ([]Node, error) {
	if n.IsElement() || n.IsRoot() {
		result := make([]Node, 0)
		for _, child := range n.children {
			if child.IsElement() && strings.EqualFold(child.GetName(), name) {
				result = append(result, child)
			}
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) GetName() string {
	return n.start.GetValues()[0]
}
func (n *Node) SetName(name string) {
	n.start.SetValue(0, name)
}
func (n *Node) getAttributes() []string {
	return n.start.GetValues()[1:]
}
func (n *Node) setAttributes(values []string) {
	for i, v := range values {
		n.start.SetValue(i+1, v)
	}
}
func (n *Node) IsNil(i int) bool {
	return n.start.IsNil(i + 1)
}
func (n *Node) SetNil(i int) {
	n.start.SetNil(i + 1)
}
func (n *Node) UnsetNil(i int) {
	n.start.UnsetNil(i + 1)
}
func (n *Node) GetComment() string {
	return n.start.GetComment()
}
func (n *Node) SetComment(comment string) error {
	return n.start.SetComment(comment)
}
func (n *Node) GetEndComment() string {
	return n.end.GetComment()
}
func (n *Node) SetEndComment(comment string) error {
	if n.IsElement() {
		return n.end.SetComment(comment)
	} else {
		return fmt.Errorf("Not an element")
	}
}
func (n *Node) AddElement(name string) (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		start, _ := wsv.NewLine([]string{name}, nil, "")
		end, _ := wsv.NewLine([]string{"end"}, nil, "")
		node := Node{start: &start, end: &end, children: make([]Node, 0)}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) AddAttribute(name string, values []string) (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		v := make([]string, len(values)+1)
		v = append(v, name)
		v = append(v, values...)
		start, _ := wsv.NewLine(v, nil, "")
		node := Node{start: &start}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) AddEmpty() (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		start, _ := wsv.NewLine(nil, nil, "")
		node := Node{start: &start}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}

func (n *Node) GetSpaces() []string {
	return n.start.GetSpaces()
}
func (n *Node) SetSpaces(spaces []string) error {
	return n.start.SetSpaces(spaces)
}
func (n *Node) GetEndSpaces() []string {
	return n.end.GetSpaces()
}
func (n *Node) SetEndSpaces(spaces []string) error {
	return n.end.SetSpaces(spaces)
}
func (n *Node) Each(fn func(n *Node) error) error {
	for _, child := range n.children {
		if err := fn(&child); err != nil {
			return err
		}
	}
	return nil
}
func (n *Node) Filter(fn func(n *Node) bool) []Node {
	result := make([]Node, 0)
	for _, child := range n.children {
		if fn(&child) {
			result = append(result, child)
		}
	}
	return result
}

func (n *Node) Write(w io.Writer) error {
	io.WriteString(w, n.start.String())
	if err := n.Each(func(n *Node) error {
		return n.Write(w)
	}); err != nil {
		return err
	}
	io.WriteString(w, n.end.String())

	return nil
}

func (n *Node) AlignAttributes(spacesBetween string, maxColumns int, rightAligned []bool) error {
	if err := wsv.ValidateSpace(spacesBetween, false); err != nil {
		return err
	}

	attributes := n.Filter(func(n *Node) bool { return n.IsAttribute() })
	spacesArray := make([][]string, 0)
	var columns int
	for _, a := range attributes {
		spacesArray = append(spacesArray, make([]string, 0))
		c := len(a.start.GetValues())
		if c > columns {
			columns = c
		}
	}

	if maxColumns > columns {
		columns = maxColumns
	}
	for i := 0; i < columns; i++ {
		var maxLen int
		for j, a := range attributes {
			values := a.start.GetValues()
			if j >= len(values) {
				continue
			}
			value := values[j]

			serValue := wsv.SerializeValue(value, false)
			cpLen := utf8.RuneCountInString(serValue)
			if cpLen > maxLen {
				maxLen = cpLen
			}
		}
		for j, a := range attributes {
			values := a.start.GetValues()
			if j >= len(values) {
				continue
			}
			value := values[j]
			serValue := wsv.SerializeValue(value, false)
			lengthDif := maxLen - utf8.RuneCountInString(serValue)
			fillingWhitespace := strings.Repeat(" ", lengthDif)
			if rightAligned[j] {
				last := spacesArray[j][len(spacesArray[j])-1]
				spacesArray[j][len(spacesArray[j])-1] = last + fillingWhitespace
				if i >= len(values)-1 {
					continue
				}
				spacesArray[j] = append(spacesArray[j], spacesBetween)
			} else {
				if i >= len(values)-1 {
					continue
				}
				spacesArray[j] = append(spacesArray[j], fillingWhitespace+spacesBetween)
			}
		}
	}
	for i, a := range attributes {
		a.SetSpaces(spacesArray[i])
	}
	return nil
}
