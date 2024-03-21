package sml

type stack []*Node

func (s *stack) Push(n *Node) {
	*s = append(*s, n)
}

func (s *stack) Pop() *Node {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}
