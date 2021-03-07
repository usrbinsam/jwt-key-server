package routes

import (
	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")

	users.POST("/register", func(context *gin.Context) {

	})

	users.POST("/login", func(context *gin.Context) {

	})
}
