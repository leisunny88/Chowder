package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AllRouters() {
	r := gin.Default()
	v1Group := r.Group("/api/v1")
	{
		IndexRouter(v1Group)
		InitDBRouter(v1Group)
		RegisterRouter(v1Group)
		LoginRouter(v1Group)
		DeviceRouter(v1Group)
		ServerDeviceRouter(v1Group)
		WeatherDataRouter(v1Group)
	}
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
