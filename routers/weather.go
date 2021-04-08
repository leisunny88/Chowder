package routers

import (
	"gin-project/api"
	"github.com/gin-gonic/gin"
)

func WeatherDataRouter(e *gin.RouterGroup) {
	weather := e.Group("/weather")
	{
		weather.GET("/city/:name", api.GetWeatherData)
	}
}
