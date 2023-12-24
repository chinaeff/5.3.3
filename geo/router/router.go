package router

import (
	"geotask_pprof/auth"
	"geotask_pprof/geo/module/courierfacade/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	courier *controller.CourierController
}

func NewRouter(courier *controller.CourierController) *Router {
	return &Router{courier: courier}
}

func (r *Router) Swagger(router *gin.RouterGroup) {
	router.GET("/swagger", swaggerUI)
}

func (r *Router) CourierAPI(router *gin.RouterGroup) {
	securedRoutes := router.Group("/secured")
	securedRoutes.Use(auth.AuthMiddleware())
	{
		securedRoutes.POST("/some-protected-endpoint", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Authorized access"})
		})
	}
}
