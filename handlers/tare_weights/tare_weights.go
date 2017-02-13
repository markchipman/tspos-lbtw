package tare_weights

import (
	"github.com/wormling/tspos-lbtw/models"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func Create(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	tareWeight := models.TareWeight{}
	err := c.BindJSON(&tareWeight)
	if err != nil {
		c.Error(err)
		return
	}

	err = db.C(models.CollectionTareWeights).Insert(tareWeight)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, tareWeight)
}

func Get(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	tareWeight := models.TareWeight{}
	oID := bson.ObjectIdHex(c.Param("_id"))
	err := db.C(models.CollectionTareWeights).FindId(oID).One(&tareWeight)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, tareWeight)
}

func List(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	tareWeights := []models.TareWeight{}
	err := db.C(models.CollectionTareWeights).Find(nil).Sort("-updated_on").All(&tareWeights)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"tareWeights": tareWeights,
	})
}

func Update(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	tareWeight := models.TareWeight{}
	err := c.Bind(&tareWeight)
	if err != nil {
		c.Error(err)
		return
	}

	query := bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))}
	doc := bson.M{
		"brand":      tareWeight.Brand,
		"category":   tareWeight.Category,
		"name":       tareWeight.Name,
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = db.C(models.CollectionTareWeights).Update(query, doc)
	if err != nil {
		c.Error(err)
	}
}

func Delete(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	query := bson.M{"_id": bson.ObjectIdHex(c.Param("_id"))}
	err := db.C(models.CollectionTareWeights).Remove(query)
	if err != nil {
		c.Error(err)
	}
}
