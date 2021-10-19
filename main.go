package main

import "github.com/josexy/dstartlite-gosdl/dstarlite"

/*
运行后
Esc/Q键退出程序，空格键单步执行算法
鼠标左键双击设置起点，右键双击设置终点
按住鼠标左键移动添加障碍物，按住鼠标中建移动移除障碍物
have fun~
*/
func main() {
	window := dstarlite.NewWindow(500, 500, "DStarLite")
	window.Run()
}
