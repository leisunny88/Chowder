package main

import (
	_ "gin-project/core"
	"gin-project/global"
	"gin-project/initialize"
	"gin-project/routers"
)

func main() {
	global.DB = initialize.InitDB() // 初始化数据库
	routers.AllRouters()            // 加载应用路由
}
