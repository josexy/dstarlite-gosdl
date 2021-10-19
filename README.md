# dstarlite-gosdl

Go implement D*Lite with SDL2

使用Go SDL2库实现的D*Lite算法

注意，尽可能使用 go mod ，自动下载依赖库
```shell
git clone https://github.com/josexy/dstarlite-gosdl.git
cd dstarlite-gosdl
go run main.go
```

运行后

- Esc/Q键退出程序，空格键单步执行算法

- 鼠标左键双击设置起点，右键双击设置终点

- 按住鼠标左键移动添加障碍物，按住鼠标中建移动移除障碍物

- have fun~

程序运行截图
![run](https://github.com/josexy/dstarlite-gosdl/blob/main/run.jpg)

具体的D*Lite论文地址: https://www.aaai.org/Papers/AAAI/2002/AAAI02-072.pdf
