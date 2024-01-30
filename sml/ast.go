package sml

import (
	"fmt"
	"math/big"
	"strings"
	"wsv"
)

type Kind int

const (
	Root      Kind = 0
	Element        = 1
	Attribute      = 2
	Empty          = 3
)

type Node struct {
	kind       Kind
	comment    string
	endComment string
	name       string
	values     []string
	spaces     []string
	endSpaces  []string
	nulls      big.Int
	children   []Node
}

func NewRoot() *Node {
	return &Node{kind: Root}
}

func (n *Node) IsRoot() bool {
	return n.kind == Root
}
func (n *Node) IsElement() bool {
	return n.kind == Element
}
func (n *Node) IsAttribute() bool {
	return n.kind == Attribute
}
func (n *Node) IsEmpty() bool {
	return n.kind == Empty
}
func (n *Node) FilterAttributes(name string) []Node {
	result := make([]Node, 0)
	for i := 0; i < len(n.children); i++ {
		n := n.children[i]
		if n.IsAttribute() && n.name == name {
			result = append(result, n)
		}
	}
	return result
}
func (n *Node) FilterElements(name string) []Node {
	result := make([]Node, 0)
	for i := 0; i < len(n.children); i++ {
		n := n.children[i]
		if n.IsElement() && n.name == name {
			result = append(result, n)
		}
	}
	return result
}
func (n *Node) GetName() string {
	return n.name
}
func (n *Node) SetName(name string) error {
	n.name = name
}
func (n *Node) getValues() []string {
	return n.values
}
func (n *Node) setValues(values []string) {
	n.values = values
}
func (n *Node) IsNil(i int) bool {
	return n.nulls.Bit(i) == 1
}
func (n *Node) SetNil(i int) {
	n.nulls.SetBit(&n.nulls, i, 1)
}
func (n *Node) UnsetNil(i int) {
	n.nulls.SetBit(&n.nulls, i, 0)
}
func (n *Node) GetComment() string {
	return n.comment
}
func (n *Node) SetComment(comment string) error {
	if e := validateComment(comment); e != nil {
		return e
	}
	n.comment = comment
	return nil
}
func (n *Node) GetEndComment() string {
	return n.endComment
}
func (n *Node) SetEndComment(comment string) error {
	if n.IsElement() {
		if e := validateComment(comment); e != nil {
			return e
		}
		n.comment = comment
		return nil
	} else {
		return fmt.Errorf("Not an element")
	}
}
func (n *Node) AddElement(name string) (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		node := Node{kind: Element, name: name}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) AddAttribute(name string, values []string) (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		node := Node{kind: Attribute, name: name, values: values}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}
func (n *Node) AddEmpty() (*Node, error) {
	if n.IsElement() || n.IsRoot() {
		node := Node{kind: Empty}
		n.children = append(n.children, node)
		return &node, nil
	} else {
		return nil, fmt.Errorf("Not an element")
	}
}

func (n *Node) GetSpaces() []string {
	return n.spaces
}
func (n *Node) SetSpaces(spaces []string) error {
	if err := wsv.validateSpaces(spaces); err != nil {
		return err
	}
	n.spaces = spaces
	return nil
}
func (n *Node) GetEndSpaces() []string {
	return n.endSpaces
}
func (n *Node) SetEndSpaces(spaces []string) error {
	if err := validateSpaces(spaces); err != nil {
		return err
	}
	n.endSpaces = spaces
	return nil
}
func (n *Node) Each(fn func(n *Node) error) error {
}
func (n *Node) Map(fn func(n *Node) (Node, error)) (Node, error) {
}
func (n *Node) Filter(fn func(n *Node) bool) []Node {
	result := make([]Node, 0)
	for i := 0; i < len(n.children); i++ {
		n := n.children[i]
		if fn(&n) {
			result = append(result, n)
		}
	}
	return result
}
func (n *Node) AlignAttributes(spacesBetween string, maxColumns int, rightAligned []bool) error {
	if err := validateSpace(spacesBetween, false); err != nil {
		return err
	}

	attributes := n.Filter(IsAttribute)
	spacesArray := make([][]string, 0)
	valuesArray := make([][]string, 0)
	var columns int
	for i := 0; i < len(attributes); i++ {
		a := attributes[i]
		spacesArray = append(spacesArray, make([]string))
		values := make([]string)
		values = append(values, a.name)
		values = append(values, a.values...)
		valuesArray = append(valuesArray, values)
		c := len(values)
		if c > columns {
			columns = c
		}
	}

	if maxColumns > columns {
		columns = maxColumns
	}
	for i := 0; i < columns; i++ {
		var maxLen int
		for j := 0; j < len(attributes); j++ {
			values := valuesArray[j]
			if i >= len(values) {
				continue
			}
			value := values[i]
			serValue := wsv.Serialize(value)
			cpLen := wsv.CountCodepoints(serValue)
			if cpLen > maxLen {
				maxLen = cpLen
			}
		}
		for j := 0; j < len(attributes); j++ {
			values := valuesArray[j]
			if i >= len(values) {
				continue
			}
			value := values[j]
			serValue := wsv.Serialize(value)
			lengthDif := maxLen - wsv.CountCodePoints(serValue)
			fillingWhitespace := strings.Repeat(" ", lengthDif)
			if rightAligned[i] { // TODO could be a panic here
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
	for i := 0; i < len(attributes); i++ {
		attributes[i].spaces = spacesArray[i]

	}
}

func validateComment(s string) error {
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		cp := r[i]
		if cp == 0x000A {
			return fmt.Errorf("Line feed in comment is not allowed")
		} else if cp >= 0xD800 && cp <= 0xDFFF {
			i++
			if cp >= 0xDC00 || i >= len(s) {
				return fmt.Errorf("Invalid UTF-16 String")
			}
			cp2 := r[i]
			if !(cp2 >= 0xDC00 && cp2 <= 0xDFFF) {
				return fmt.Errorf("Invalid UTF-16 String")
			}
		}
	}
	return nil
}
