package tare_weights

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-contrib/location"
	"github.com/wormling/tspos-lbtw/models"
)

func Create(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	tareWeight := models.TareWeight{}
	err := c.BindJSON(&tareWeight)
	if err != nil {
		c.Error(err)
		return
	}

	tareWeight.CreatedOn = time.Now().UnixNano() / int64(time.Millisecond)
	tareWeight.UpdatedOn = time.Now().UnixNano() / int64(time.Millisecond)

	err = db.C(models.CollectionTareWeights).Insert(tareWeight)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.JSON(http.StatusOK, tareWeight)
	}
}

func List(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	tareWeights := []models.TareWeight{}
	//query := models.TareWeight{}

	var query map[string]string
	var z []string
	var ss []string
	var page int = 0
	var per_page int = 10

	ss = strings.Split(c.Request.URL.RawQuery, "&")
	query = make(map[string]string)
	for _, pair := range ss {
		z = strings.Split(pair, "=")
		if len(z) > 1 {
			switch z[0] {
			case "page":
				page, _ = strconv.Atoi(z[1])
			case "per_page":
				per_page, _ = strconv.Atoi(z[1])
			default:
				query[z[0]], _ = url.QueryUnescape(z[1])
			}

		}
	}

	err := db.C(models.CollectionTareWeights).Find(query).Skip(page * per_page).Limit(per_page).Sort("-updated_on").All(&tareWeights)
	if err != nil {
		c.Error(err)
	}

	var count int = 0
	count, err = db.C(models.CollectionTareWeights).Find(query).Count()

	// rfc5988
	links := MakeLinkHeader(c, page, per_page, count)
	c.Header("Link", links)
	c.Header("X-Total-Count", string(count))

	c.JSON(http.StatusOK, tareWeights)
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

	tareWeight.UpdatedOn = time.Now().UnixNano() / int64(time.Millisecond)

	err = db.C(models.CollectionTareWeights).Update(query, tareWeight)
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

func MakeLink(c *gin.Context, page int, perPage int, rel_name string) string {
	lurl := location.Get(c)
	query := c.Request.URL.Query()
	query.Set("page", string(page))
	query.Set("per_page", string(perPage))
	link := c.Request
	link.URL.RawQuery = query.Encode()

	return fmt.Sprintf("<%s>; rel=\"%s\" ", lurl.Scheme+"://"+lurl.Host+link.RequestURI, rel_name)
}

func MakeLinkHeader(c *gin.Context, page int, perPage int, count int) string {
	s := ""

	// Build first link
	s += MakeLink(c, 0, perPage, "first")

	// Build last link
	s += MakeLink(c, count/perPage, perPage, "last")

	if page >= 1 {
		// Build prev link
		s += MakeLink(c, page-1, perPage, "prev")
	}

	if page <= count {
		// Build next link
		s += MakeLink(c, page+1, perPage, "next")
	}

	return s
}
