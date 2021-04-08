package routers

import (
	"gin-project/api"
	"github.com/gin-gonic/gin"
)

func IndexRouter(e *gin.RouterGroup) {
	e.GET("/index", api.Index)
}


//func IndexRouter(e *gin.RouterGroup) {
//	e.GET("/index", api.Index)
//}