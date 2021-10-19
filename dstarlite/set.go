package dstarlite

type void struct{}

var member void

type Set struct {
	set map[*Node]void
}

func NewSet() *Set {
	return &Set{
		set: make(map[*Node]void),
	}
}

func (s *Set) Clear() {
	s.set = make(map[*Node]void)
}

func (s *Set) Range(fn func(node *Node)) {
	for k := range s.set {
		fn(k)
	}
}

func (s *Set) Add(node *Node) {
	s.set[node] = member
}

func (s *Set) Remove(node *Node) {
	delete(s.set, node)
}

func (s *Set) Size() int {
	return len(s.set)
}
