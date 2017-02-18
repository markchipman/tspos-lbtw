package tare_weights

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-contrib/location"
	"github.com/wormling/tspos-lbtw/models"
	"net/url"
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
	lurl := location.Get(c)
	db := c.MustGet("db").(*mgo.Database)
	tareWeights := []models.TareWeight{}
	//query := models.TareWeight{}

	var query map[string]string
	var ss []string
	var page int = 0
	var per_page int = 10

	ss = strings.Split(c.Request.URL.RawQuery, "&")
	query = make(map[string]string)
	for _, pair := range ss {
		z := strings.Split(pair, "=")
		switch z[0] {
		case "page":
			page, _ = strconv.Atoi(z[1])
		case "per_page":
			per_page, _ = strconv.Atoi(z[1])
		default:
			query[z[0]], _ = url.QueryUnescape(z[1])
		}

	}

	err := db.C(models.CollectionTareWeights).Find(query).Skip(page * per_page).Limit(per_page).Sort("-updated_on").All(&tareWeights)
	if err != nil {
		c.Error(err)
	}

	var count int = 0
	count, err = db.C(models.CollectionTareWeights).Find(query).Count()

	// Build first link
	firstQuery := c.Request.URL.Query()
	firstQuery.Set("page", "0")
	firstQuery.Set("per_page", string(per_page))
	first := c.Request
	first.URL.RawQuery = firstQuery.Encode()
	firstURL := lurl.Scheme + "://" + lurl.Host + first.RequestURI

	// Build last link
	lastQuery := c.Request.URL.Query()
	lastQuery.Set("page", string(count/per_page))
	lastQuery.Set("per_page", string(per_page))
	last := c.Request
	last.URL.RawQuery = lastQuery.Encode()
	lastURL := lurl.Scheme + "://" + lurl.Host + last.RequestURI

	// Build prev link
	prevQuery := c.Request.URL.Query()
	prevQuery.Set("page", string(page))
	prevQuery.Set("per_page", string(per_page))
	prev := c.Request
	prev.URL.RawQuery = prevQuery.Encode()
	prevURL := lurl.Scheme + "://" + lurl.Host + prev.RequestURI

	// Build next link
	nextQuery := c.Request.URL.Query()
	nextQuery.Set("page", string(page))
	nextQuery.Set("per_page", string(per_page))
	next := c.Request
	next.URL.RawQuery = nextQuery.Encode()
	nextURL := lurl.Scheme + "://" + lurl.Host + next.RequestURI

	// rfc5988
	links := fmt.Sprintf("<%s>; rel=\"next\" <%s>; rel=\"prev\" <%s>; rel=\"first\" <%s>; rel=\"last\"", prevURL, nextURL, firstURL, lastURL)
	c.Header("Link", links)

	c.JSON(http.StatusOK, gin.H{
		"tareWeights": tareWeights,
		"hmm":         lurl.Scheme,
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
