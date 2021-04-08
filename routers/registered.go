package routers

import "github.com/gin-gonic/gin"
import "gin-project/api"

func RegisterRouter(e *gin.RouterGroup) {
	userRouter := e.Group("/user")
	{
		userRouter.POST("/register", api.Register)
	}
	//e.POST("/user/register", api.Register)
}
