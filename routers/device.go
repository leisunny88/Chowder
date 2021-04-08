package routers

import (
	"gin-project/api"
	"github.com/gin-gonic/gin"
)

func DeviceRouter(e *gin.RouterGroup) {
	device := e.Group("/device")
	{
		device.POST("/create", api.BatchCreateDevice)
		device.DELETE("/delete", api.BatchDeleteDevice)
		device.GET("/mem", api.GetDeviceMem)
		device.GET("/cpu", api.GetDeviceCPU)
		device.GET("/disk/:id", api.GetDeviceDisk)
		device.GET("/host/:id", api.GetHostInfo)
		device.GET("/process/:id", api.GetDeviceProcess)
	}
	//e.POST("/create/device", api.BatchCreateDevice)
	//e.DELETE("/delete/device/:id", api.DeleteDevice)
}
