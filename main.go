package main

import (
	"net/http"

	"github.com/gin-contrib/location"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"

	"github.com/wormling/tspos-lbtw/config"
	"github.com/wormling/tspos-lbtw/db"
	"github.com/wormling/tspos-lbtw/handlers/tare_weights"
	"github.com/wormling/tspos-lbtw/middlewares"
)

func init() {
	db.Connect()
}

func main() {
	config.BuildDefaultConf()
	c, _ := config.LoadConfYaml("./config.yaml")

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

	go ensureIndex()

	// Start Server
	bind := c.Core.Listener.Bind
	port := c.Core.Listener.Port
	router.Run(bind + ":" + port)
}

func ensureIndex() {
	session := db.Session.Clone()
	defer session.Close()

	c := session.DB("tspos_lbtw").C("tare_weights")

	index := mgo.Index{
		//Key:        []string{"brand", "category", "name", "bottle_size, empty_weight, full_weight, image_url, created_on, updated_on"},
		Key:        []string{"brand"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Collation: &mgo.Collation{
			Locale:   "en",
			Strength: 2,
		},
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
