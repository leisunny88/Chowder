package routers

import (
	"gin-project/api"
	"github.com/gin-gonic/gin"
)

func LoginRouter(e *gin.RouterGroup) {
	userRouter := e.Group("/user")
	{
		userRouter.POST("/login", api.Login)
	}
	//e.POST("/user/login", api.Login)
}
