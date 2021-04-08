package routers

import (
	"gin-project/api"
	"github.com/gin-gonic/gin"
)

func ServerDeviceRouter(e *gin.RouterGroup) {
	DevSSh := e.Group("/server-device")
	{
		DevSSh.GET("/cmd-data", api.ServerDevice)
		DevSSh.GET("/terminal", api.SimulationTerminal)
	}
}
