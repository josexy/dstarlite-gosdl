package dstarlite

import (
	"fmt"
	"math"
)

const Inf = math.MaxFloat64

type Key struct {
	k1, k2 float64
}

func (k Key) Less(o Key) bool {
	if k.k1 == o.k1 {
		return k.k2 < o.k2
	}
	return k.k1 < o.k1
}

type Node struct {
	x, y      int
	g, h, rhs float64
	k         Key
	obstacle  bool
}

func NewNode(x, y int) *Node {
	return &Node{
		g: Inf, h: Inf, rhs: Inf,
		x: x, y: y,
		obstacle: false,
	}
}

func (n Node) Equal(o *Node) bool {
	return n.x == o.x && n.y == o.y
}

func (n Node) String() string {
	return fmt.Sprintf("(%d, %d)", n.x, n.y)
}
