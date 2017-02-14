package main

import (
	"flag"
	"fmt"
	"github.com/wormling/tspos-lbtw/db"
	"github.com/wormling/tspos-lbtw/handlers/tare_weights"
	"github.com/wormling/tspos-lbtw/middlewares"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Config Configuration information
type Config struct {
	Listener struct {
		Bind string `yaml:"bind"`
		Port string `yaml:"port"`
	}
	Database struct {
		Url string `yaml:"url"`
	}
}

func init() {
	db.Connect()
}

func main() {
	// Get Arguments
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "./config.yaml", "Path to Config File")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command> \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	config := Config{}
	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatal("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true

	router.Use(middlewares.Connect)
	router.Use(middlewares.ErrorHandler)

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

	ensureIndex()

	// Start Server
	bind := config.Listener.Bind
	port := config.Listener.Port
	router.Run(bind + ":" + port)
}

func ensureIndex() {
	session := db.Session.Clone()
	defer session.Close()

	c := session.DB("tspos_lbtw").C("tare_weights")

	index := mgo.Index{
		Key:        []string{"brand", "category", "name", "bottle_size"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
