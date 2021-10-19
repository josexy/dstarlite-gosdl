package dstarlite

var offset = [8][2]int{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}

type Grid struct {
	Vd     [][]*Node
	width  int
	height int
}

func NewGrid(w, h int) *Grid {
	return &Grid{width: w, height: h}
}

func (g *Grid) reset() {
	for x := 0; x < g.height; x++ {
		for y := 0; y < g.width; y++ {
			node := g.Vd[x][y]
			node.k.k1, node.k.k2 = 0, 0
			node.g, node.h, node.rhs = Inf, Inf, Inf
		}
	}
}

func (g *Grid) cell(x, y int) *Node {
	if g.in(x, y) {
		return g.Vd[x][y]
	}
	return nil
}

func (g *Grid) in(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.height && y < g.width
}

func (g *Grid) GetPredecessors(node *Node) []*Node {
	return g.GetSuccessors(node)
}

func (g *Grid) GetSuccessors(node *Node) (list []*Node) {
	for i := 0; i < 8; i++ {
		nx, ny := node.x+offset[i][0], node.y+offset[i][1]
		if g.in(nx, ny) && !g.Vd[nx][ny].obstacle {
			list = append(list, g.Vd[nx][ny])
		}
	}
	return
}
