package main

import (
	"github.com/wormling/tspos-lbtw/config"
	"github.com/wormling/tspos-lbtw/db"
	"github.com/wormling/tspos-lbtw/routes"
)

var conf config.ConfYaml

func init() {
	config.BuildDefaultConf()
	conf, _ = config.LoadConfYaml("./config.yaml")
	config.BuildDefaultConf()
	db.Connect(conf.Core.Database.Url)
	db.EnsureIndex()
}

func main() {
	router := routes.NewRouter()

	// Start Server
	bind := conf.Core.Listener.Bind
	port := conf.Core.Listener.Port
	router.Run(bind + ":" + port)
}
