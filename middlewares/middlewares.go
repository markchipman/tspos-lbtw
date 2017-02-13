// Package middlewares contains gin middlewares
// Usage: router.Use(middlewares.Connect)
package middlewares

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/wormling/tspos-lbtw/db"
	"net/http"
)

// Connect middleware clones the database session for each request and
// makes the `db` object available for each handler
func Connect(c *gin.Context) {
	s := db.Session.Clone()

	defer s.Close()

	c.Set("db", s.DB(db.Mongo.Database))
	c.Next()
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.HTML(http.StatusBadRequest, "400", gin.H{
			"errors": c.Errors,
		})
	}
}
