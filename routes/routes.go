package routes

import (
	"net/http"

	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/gin-contrib/location"
	"github.com/wormling/tspos-lbtw.v1/handlers"
	"github.com/wormling/tspos-lbtw.v1/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true

	// configure to automatically detect scheme and host
	// - use http when default scheme cannot be determined
	// - use localhost:8080 when default host cannot be determined
	router.Use(location.Default())
	router.Use(middlewares.Connect)
	router.Use(middlewares.ErrorHandler)
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/v1/tare/weights")
	})
	v1 := router.Group("/v1/tare")
	{
		v1.GET("/weights/:_id", tare_weights.Get)
		v1.GET("/weights", tare_weights.List)
		v1.POST("/weights", tare_weights.Create)
		v1.PUT("/weights/:_id", tare_weights.Update)
		v1.DELETE("/weights/:_id", tare_weights.Delete)
	}

	return router
}
