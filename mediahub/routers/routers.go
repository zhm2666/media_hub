package routers

import (
	"github.com/gin-gonic/gin"
	"mediahub/controller"
)

func InitRouters(api *gin.RouterGroup, c *controller.Controller) {
	v1 := api.Group("/v1")
	v1.POST("/file/upload", c.Upload)
	v1.GET("/home", c.Home)
}
