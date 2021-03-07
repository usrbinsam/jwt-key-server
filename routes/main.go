package routes

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(router *gin.Engine) {
	api := router.Group("/api")
	addUserRoutes(api)
}
