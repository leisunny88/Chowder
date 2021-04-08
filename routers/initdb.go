package routers

import (
	"gin-project/initialize"
	"github.com/gin-gonic/gin"
)

func InitDBRouter(e *gin.RouterGroup) {
	e.GET("/init-db", initialize.Migrate)
}
