package models

import "gopkg.in/mgo.v2/bson"

const (
	// CollectionTareWeights Collection name
	CollectionTareWeights = "tare_weights"
)

// TareWeight Model for liquor bottle tare weights
type TareWeight struct {
	Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Brand       string        `json:"brand" form:"brand" binding:"required" bson:"brand"`
	Category    string        `json:"category" form:"category" binding:"required" bson:"category"` // Rum, Vodka, Gin, ...
	Name        string        `json:"name" form:"name" binding:"required" bson:"name"`
	BottleSize  float32       `json:"bottle_size" form:"bottle_size" binding:"required" bson:"bottle_size"`    // Size in ml
	EmptyWeight float32       `json:"empty_weight" form:"empty_weight" binding:"required" bson:"empty_weight"` // Tare weight in grams
	FullWeight  float32       `json:"full_weight" form:"full_weight" binding:"required" bson:"full_weight"`    // Full weight in grams
	ImageUrl    string        `json:"image_url" form:"image_url" bson:"image_url"`
	CreatedOn   int64         `json:"created_on" bson:"created_on"`
	UpdatedOn   int64         `json:"updated_on" bson:"updated_on"`
}
