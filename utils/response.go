package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseCodeMsg(c *gin.Context, m string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": m,
	})
}

func ResponseNotFoundCode(c *gin.Context, m string) {
	c.JSON(http.StatusNotFound, gin.H{
		"code": http.StatusNotFound,
		"msg": m,
	})
}

func ResponseJsonMsg(c *gin.Context, m interface{})  {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": m,
	})
}

func ResponseCreateMsg(c *gin.Context, m string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusCreated,
		"msg": m,
	})
}