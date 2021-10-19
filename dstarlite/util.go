package dstarlite

import "math"

func heuristic(n1 *Node, n2 *Node) float64 {
	dx, dy := n1.x-n2.x, n1.y-n2.y
	// 曼哈顿距离
	return 7 * (math.Abs(float64(dx)) + math.Abs(float64(dy)))
	// 欧几里得距离
	// return 7 * math.Sqrt(float64(dx*dx+dy*dy))
}

func cost(n1 *Node, n2 *Node) float64 {
	if n1.obstacle || n2.obstacle {
		return Inf
	}
	if n1.x-n2.x == 0 || n1.y-n2.y == 0 {
		return 10
	}
	return 14
}
