package api

import (
	"gin-project/utils"
	"github.com/gin-gonic/gin"
)

// 服务器设备连接
func ServerDevice(c *gin.Context) {
	utils.SSHClient(c)
}

// 模拟交互终端
func SimulationTerminal(c *gin.Context) {
	utils.Terminal(c)
}
