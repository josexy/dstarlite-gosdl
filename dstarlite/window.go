package dstarlite

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const CellSize = 10

type Window struct {
	U                    *PriorityQueue // 优先级队列
	S                    *Set           // 增加或删除的障碍物集合
	C                    *Set           // 计算轨迹
	grid                 *Grid          // 栅格图
	path                 []*Node        // 移动轨迹
	window               *sdl.Window
	render               *sdl.Renderer
	width, height        int
	title                string
	running              bool
	km                   float64
	sStart, sGoal, sLast *Node
	startX, startY       int
	goalX, goalY         int
	isStart              bool
	isChangeObstacle     bool
	hasObstacle          bool
}

func NewWindow(width, height int, title string) *Window {
	w := &Window{
		width: width, height: height, title: title,
		running: true,
		U:       NewPriorityQueue(),
		S:       NewSet(), C: NewSet(),
		grid: NewGrid(width/CellSize, height/CellSize),
	}
	w.init()
	return w
}

func (w *Window) init() {
	w.goalX = w.grid.height - 1
	w.goalY = w.grid.width - 1
	w.initMap2D()

	w.sStart = w.grid.cell(w.startX, w.startY)
	w.sGoal = w.grid.cell(w.goalX, w.goalY)
}

func (w *Window) Run() {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	w.window, err = sdl.CreateWindow(w.title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(w.width), int32(w.height), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	w.render, err = sdl.CreateRenderer(w.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	for w.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				w.running = false
				break
			case *sdl.KeyboardEvent:
				w.keyPressEvent(event.(*sdl.KeyboardEvent))
			case *sdl.MouseMotionEvent:
				w.mouseMoveEvent(event.(*sdl.MouseMotionEvent))
			case *sdl.MouseButtonEvent:
				w.mouseButtonEvent(event.(*sdl.MouseButtonEvent))
			}
		}

		w.update()
		sdl.Delay(16)
	}
	w.Quit()
}

func (w *Window) Quit() {
	_ = w.render.Destroy()
	_ = w.window.Destroy()
	sdl.Quit()
}

func (w *Window) keyPressEvent(ke *sdl.KeyboardEvent) {
	switch ke.Keysym.Sym {
	case sdl.K_ESCAPE, sdl.K_q:
		w.running = false
	case sdl.K_SPACE:
		if ke.State == sdl.PRESSED {
			if !w.isStart {
				w.run()
			}
			w.tick()
		}
	}
}

func (w *Window) mouseMoveEvent(me *sdl.MouseMotionEvent) {
	x, y := mapToGridPoint(me.X, me.Y)
	if me.State == sdl.BUTTON_MIDDLE {
		node := w.grid.cell(x, y)
		if node == nil {
			return
		}
		// 等待释放按键时触发重新计算当前局势的最短路径
		w.isChangeObstacle = true
		// 障碍物删除
		node.obstacle = false
		w.S.Add(node)
		w.update()
	} else if me.State == sdl.BUTTON_LEFT {
		node := w.grid.cell(x, y)
		if node == nil {
			return
		}
		// 等待释放按键时触发重新计算当前局势的最短路径
		w.isChangeObstacle = true
		node.obstacle = true
		tmpS := NewSet()
		// 绘制障碍物
		rangeNum := 2
		for _x := -rangeNum; _x < rangeNum; _x++ {
			for _y := -rangeNum; _y < rangeNum; _y++ {
				t := w.grid.cell(x+_x, y+_y)
				if t != nil {
					// 障碍物增加
					t.obstacle = true
					tmpS.Add(t)
				}
			}
		}
		tmpS.Range(func(node *Node) {
			succ := w.grid.GetSuccessors(node)
			for _, s := range succ {
				w.S.Add(s)
			}
		})
		w.update()
	}
}

func (w *Window) mouseButtonEvent(be *sdl.MouseButtonEvent) {
	// 释放鼠标，此时触发重新计算当前局势的最短路径
	if be.State == sdl.RELEASED {
		if w.isChangeObstacle {
			// 防止重复触发
			w.isChangeObstacle = false
			w.hasObstacle = true
		}
	} else if !w.isStart && be.Clicks == 2 && be.State == sdl.PRESSED {
		node := w.grid.cell(mapToGridPoint(be.X, be.Y))
		if be.Button == sdl.BUTTON_LEFT {
			w.sStart = node
			w.startX, w.startY = node.x, node.y
			fmt.Printf("重新设置起点: %v\n", w.sStart)
		} else if be.Button == sdl.BUTTON_RIGHT {
			w.sGoal = node
			w.goalX, w.goalY = node.x, node.y
			fmt.Printf("重新设置终点: %v\n", w.sStart)
		}
	}
	w.update()
}

// mapToScreenRect 将二维坐标映射到窗口坐标
func mapToScreenRect(x, y int) sdl.Rect {
	return sdl.Rect{X: int32(y) * CellSize, Y: int32(x) * CellSize, W: CellSize, H: CellSize}
}

// mapToGridPoint 将窗口坐标映射到二维坐标
func mapToGridPoint(x, y int32) (int, int) {
	return int(y) / CellSize, int(x) / CellSize
}

func (w *Window) initMap2D() {
	for x := 0; x < w.grid.height; x++ {
		var l []*Node
		for y := 0; y < w.grid.width; y++ {
			l = append(l, NewNode(x, y))
		}
		w.grid.Vd = append(w.grid.Vd, l)
	}
}

func (w *Window) clear() {
	// 背景色白色
	_ = w.render.SetDrawColor(255, 255, 255, 0)
	_ = w.render.Clear()
}

func (w *Window) drawRect(r, g, b uint8, rect sdl.Rect) {
	_ = w.render.SetDrawColor(r, g, b, 0)
	_ = w.render.DrawRect(&rect)
}

func (w *Window) fillRect(r, g, b uint8, rect sdl.Rect) {
	_ = w.render.SetDrawColor(r, g, b, 0)
	_ = w.render.FillRect(&rect)
}

func (w *Window) update() {
	w.clear()

	// 绘制背景
	for x := 0; x < w.grid.height; x++ {
		for y := 0; y < w.grid.width; y++ {
			rect := mapToScreenRect(x, y)
			node := w.grid.cell(x, y)
			if x == w.startX && y == w.startY {
				w.fillRect(0, 255, 0, rect) // 起点，绿色
			} else if x == w.goalX && y == w.goalY {
				w.fillRect(255, 0, 0, rect) // 终点，红色
			} else if node.obstacle {
				w.fillRect(100, 100, 100, rect) // 障碍物，灰色
			} else {
				w.drawRect(0, 0, 0, rect)
			}
		}
	}

	// 绘制计算节点
	w.C.Range(func(node *Node) {
		w.drawRect(200, 200, 10, mapToScreenRect(node.x, node.y))
	})

	// 绘制行走路径
	for _, s := range w.path {
		w.fillRect(0, 255, 0, mapToScreenRect(s.x, s.y))
	}

	// 绘制当前正在移动的节点
	if w.sStart != nil {
		w.fillRect(0, 0, 255, mapToScreenRect(w.sStart.x, w.sStart.y))
	}

	w.render.Present()
}

// showMessageBox 显示消息框
func (w *Window) showMessageBox(title, message string, typ uint32) {
	buttons := []sdl.MessageBoxButtonData{
		{0, 0, "OK"},
	}
	colorScheme := sdl.MessageBoxColorScheme{
		Colors: [5]sdl.MessageBoxColor{
			{255, 0, 0},
			{0, 255, 0},
			{255, 255, 0},
			{0, 0, 255},
			{255, 0, 255},
		},
	}
	mbd := sdl.MessageBoxData{
		Flags:       typ,
		Window:      w.window,
		Title:       title,
		Message:     message,
		Buttons:     buttons,
		ColorScheme: &colorScheme,
	}

	var err error
	if _, err = sdl.ShowMessageBox(&mbd); err != nil {
		fmt.Println("error displaying message box")
		return
	}
}

// run 运行 D*Lite算法
func (w *Window) run() {
	w.U.Clear()
	w.C.Clear()
	w.S.Clear()
	w.grid.reset()
	w.path = w.path[:0]
	w.update()

	w.initialize()
	w.completeShortPath()
	w.sLast = w.sStart
	w.isStart = true
}

func (w *Window) calcKey(node *Node) Key {
	val := math.Min(node.g, node.rhs)
	return Key{k1: val + heuristic(node, w.sStart) + w.km, k2: val}
}

func (w *Window) updateVertex(u *Node) {
	if u != w.sGoal {
		minRhs := Inf
		succ := w.grid.GetSuccessors(u)
		for _, s := range succ {
			minRhs = math.Min(minRhs, s.g+cost(u, s))
		}
		u.rhs = minRhs
	}
	if i := w.U.Find(u); i != -1 {
		w.U.Remove(i)
	}
	if u.g != u.rhs {
		u.k = w.calcKey(u)
		w.U.Push(u)
	}
}

func (w *Window) initialize() {
	w.km = 0
	w.sStart.h = 0
	w.sGoal.rhs = 0
	w.sGoal.k = w.calcKey(w.sGoal)
	w.U.Push(w.sGoal)
}

func (w *Window) completeShortPath() {
	w.C.Clear()
	for !w.U.IsEmpty() &&
		(w.U.Top().k.Less(w.calcKey(w.sStart)) || w.sStart.g != w.sStart.rhs) {
		u := w.U.Pop()
		w.C.Add(u)
		keyOld := u.k
		if keyOld.Less(w.calcKey(u)) {
			u.k = w.calcKey(u)
			w.U.Push(u)
		} else if u.g > u.rhs {
			u.g = u.rhs
			pred := w.grid.GetPredecessors(u)
			for _, n := range pred {
				w.updateVertex(n)
			}
		} else {
			u.g = Inf
			w.updateVertex(u)
			pred := w.grid.GetPredecessors(u)
			for _, n := range pred {
				w.updateVertex(n)
			}
		}
	}
}

// tick 单步执行算法
func (w *Window) tick() {
	if !w.isStart {
		return
	}
	if w.sStart == w.sGoal {
		fmt.Println("已经到达终点")
		w.showMessageBox("消息", "已经到达终点", sdl.MESSAGEBOX_INFORMATION)
		w.isStart = false
		w.sStart = w.grid.cell(w.startX, w.startY)
		w.sGoal = w.grid.cell(w.goalX, w.goalY)
		return
	}
	w.path = append(w.path, w.sStart)
	succ := w.grid.GetSuccessors(w.sStart)
	minCost := Inf
	// 找出下一个合适的节点
	var nextNode *Node
	for _, s := range succ {
		c := s.g + cost(w.sStart, s)
		if c < minCost {
			minCost = c
			nextNode = s
		}
	}
	if nextNode == nil {
		fmt.Println("无法移动到下一个位置")
		w.showMessageBox("警告", "无法移动到下一个位置", sdl.MESSAGEBOX_WARNING)
		w.isStart = false
		w.sStart = w.grid.cell(w.startX, w.startY)
		w.sGoal = w.grid.cell(w.goalX, w.goalY)
		return
	}

	w.update()
	fmt.Println("移动到下一个位置: ", nextNode)
	w.sStart = nextNode
	if w.hasObstacle {
		w.km += heuristic(w.sLast, w.sStart)
		w.sLast = w.sStart
		count := 0
		w.S.Range(func(node *Node) {
			// 周围环境中移除了障碍物则需要更新节点信息
			if !node.obstacle {
				w.updateVertex(node)
				count++
			}
		})
		fmt.Printf("障碍物添加或者删除数量: %d，需要更新的节点数量: %d\n", w.S.Size(), count)
		w.completeShortPath()
		w.hasObstacle = false
		w.S.Clear()
	}
}
