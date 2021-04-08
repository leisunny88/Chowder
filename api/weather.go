package api

import (
	"gin-project/service"
	"gin-project/utils"
	"github.com/gin-gonic/gin"
)

func GetWeatherData(c *gin.Context) {
	cityName := c.Param("name")
	if cityName == "" {
		service.ErrorInfo(c, "city name is null")
	}
	utils.WeatherData(c, cityName)
}