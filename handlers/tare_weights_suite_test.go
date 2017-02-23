package tare_weights_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/modocache/gory"
	"github.com/wormling/tspos-lbtw/models"
	"testing"
)

func TestTsposLbtw(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "TSPOS LBTW Suite")
}

/*
Define factories for bottle weights.
*/
func defineFactories() {
	gory.Define("tare_weight", models.TareWeight{},
		func(factory gory.Factory) {
			factory["Brand"] = "Bombay Sapphire"
			factory["Category"] = "Liquor"
			factory["Name"] = "Bombay Sapphire Gin"
			factory["BottleSize"] = 958.21
			factory["EmptyWeight"] = 688.89
			factory["FullWeight"] = 1700.0
			factory["ImageUrl"] = ""
			factory["CreatedOn"] = 0
			factory["UpdatedOn"] = 0
		})

	gory.Define("tare_weight2", models.TareWeight{},
		func(factory gory.Factory) {
			factory["Brand"] = "Bacardi"
			factory["Category"] = "Liquor"
			factory["Name"] = "Bombay Sapphire Gin"
			factory["BottleSize"] = 958.21
			factory["EmptyWeight"] = 524.74
			factory["FullWeight"] = 1600.0
			factory["ImageUrl"] = ""
			factory["CreatedOn"] = 0
			factory["UpdatedOn"] = 0
		})
}
