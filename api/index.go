package api

import (
	"gin-project/utils"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	utils.ResponseCodeMsg(c, "runing")
}
