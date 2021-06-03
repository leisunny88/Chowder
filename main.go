package main

import (
	"fmt"
	"net"
)

//import (
//	_ "gin-project/core"
//	"gin-project/global"
//	"gin-project/initialize"
//	"gin-project/routers"
//)
//
//func main() {
//	global.DB = initialize.InitDB() // 初始化数据库
//	routers.AllRouters()            // 加载应用路由
//}

func main() {
	mask := net.IPMask(net.ParseIP("192.168.100.120").To4()) // If you have the mask as a string
	//mask := net.IPv4Mask(255,255,255,0) // If you have the mask as 4 integer values

	prefixSize := mask.String()
	fmt.Println(prefixSize)
}
