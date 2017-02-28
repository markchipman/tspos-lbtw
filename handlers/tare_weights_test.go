package tare_weights_test

import (
	//. "github.com/wormling/tspos-lbtw/handlers/tare_weights"

	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"

	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"github.com/wormling/tspos-lbtw/db"
	"github.com/wormling/tspos-lbtw/models"
	"github.com/wormling/tspos-lbtw/routes"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var _ = Describe("Handlers/TareWeights", func() {
	var dbName string
	var dbUrl string
	var session *mgo.Session
	var server *gin.Engine
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up a new routes, connected to a test database,
		// before each test.
		dbName = "tspos_lbtw_test"
		dbUrl = "mongodb://localhost:27017/" + dbName
		db.Connect(dbUrl)
		session = db.Session.Clone()

		server = routes.NewRouter()

		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// Clear the database after each test.
		session.DB(dbName).DropDatabase()
		session.Close()
	})

	Describe("POST /v1/tare/weights", func() {
		var objectId = bson.NewObjectId()

		Describe("with valid input", func() {
			BeforeEach(func() {
				body, _ := json.Marshal(gory.Build("tare_weight"))
				request, _ = http.NewRequest("POST", "/v1/tare/weights", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
				collection := session.DB(dbName).C("tare_weights")
				tareWeight := gory.BuildWithParams("tare_weight", gory.Factory{
					"Id": objectId,
				}).(*models.TareWeight)
				collection.Insert(tareWeight)
			})

			Context("when tare weight is created", func() {
				It("returns a status code of 200", func() {
					fmt.Printf(recorder.Body.String())
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(200))
				})

				It("returns the created tare weight", func() {
					fmt.Printf(recorder.Body.String())
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(200))

					recorder = httptest.NewRecorder()
					request, _ = http.NewRequest("GET", "/v1/tare/weights/"+objectId.Hex(), nil)
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(200))

					var tareWeightJSON models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightJSON)

					Expect(tareWeightJSON.Brand).To(Equal("Bombay Sapphire"))
					Expect(tareWeightJSON.Category).To(Equal("Liquor"))
					Expect(tareWeightJSON.Name).To(Equal("Bombay Sapphire Gin"))
					Expect(tareWeightJSON.BottleSize).To(Equal(958.22))
					Expect(tareWeightJSON.EmptyWeight).To(Equal(688.89))
					Expect(tareWeightJSON.FullWeight).To(Equal(1700.0))
					Expect(tareWeightJSON.ImageUrl).To(Equal(""))
					Expect(tareWeightJSON.CreatedOn).To(BeTemporally("==", time.Time{}))
					Expect(tareWeightJSON.UpdatedOn).To(BeTemporally("==", time.Time{}))
				})
			})
		})

		Describe("with invalid valid input", func() {
			BeforeEach(func() {
				body, _ := json.Marshal(struct {
					Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
					Brand       string        `json:"brand" form:"brand" binding:"required" bson:"brand"`
					Category    string        `json:"category" form:"category" binding:"required" bson:"category"` // Rum, Vodka, Gin, ...
					Name        string        `json:"name" form:"name" binding:"required" bson:"name"`
					BottleSize  string        `json:"bottle_size" form:"bottle_size" binding:"required" bson:"bottle_size"`    // Size in ml
					EmptyWeight float64       `json:"empty_weight" form:"empty_weight" binding:"required" bson:"empty_weight"` // Tare weight in grams
					FullWeight  float64       `json:"full_weight" form:"full_weight" binding:"required" bson:"full_weight"`    // Full weight in grams
					ImageUrl    string        `json:"image_url" form:"image_url" bson:"image_url"`
					CreatedOn   time.Time     `json:"created_on" bson:"created_on"`
					UpdatedOn   time.Time     `json:"updated_on" bson:"updated_on"`
				}{
					bson.NewObjectId(),
					"Bombay Sapphire",
					"Liquor",
					"Bombay Sapphire Gin",
					"958.21",
					524.74,
					1600.0,
					"",
					time.Now(),
					time.Now(),
				})
				request, _ = http.NewRequest("POST", "/v1/tare/weights", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			Context("when invalid tare weight is created", func() {
				It("returns a status code of 400", func() {
					fmt.Printf(recorder.Body.String())
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(400))
				})
			})
		})
	})

	Describe("PUT /v1/tare/weights/:id", func() {
		var objectId = bson.NewObjectId()

		BeforeEach(func() {
			collection := session.DB(dbName).C("tare_weights")

			tareWeight2 := gory.BuildWithParams("tare_weight2", gory.Factory{
				"Id":        objectId,
				"CreatedOn": time.Now(),
				"UpdatedOn": time.Now(),
			}).(*models.TareWeight)
			collection.Insert(tareWeight2)

			tareWeight := gory.BuildWithParams("tare_weight", gory.Factory{
				"Id": objectId,
			}).(*models.TareWeight)

			body, _ := json.Marshal(tareWeight)
			request, _ = http.NewRequest("PUT", "/v1/tare/weights/"+tareWeight.Id.Hex(), bytes.NewReader(body))
			request.Header.Set("content-type", "application/json")
		})

		It("returns a status code of 200", func() {
			fmt.Printf(recorder.Body.String())
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
		})

		It("updated the tare weight", func() {
			fmt.Printf(recorder.Body.String())
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))

			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest("GET", "/v1/tare/weights/"+objectId.Hex(), nil)
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))

			var tareWeightJSON models.TareWeight
			json.Unmarshal(recorder.Body.Bytes(), &tareWeightJSON)

			Expect(tareWeightJSON.Brand).To(Equal("Bombay Sapphire"))
			Expect(tareWeightJSON.Category).To(Equal("Liquor"))
			Expect(tareWeightJSON.Name).To(Equal("Bombay Sapphire Gin"))
			Expect(tareWeightJSON.BottleSize).To(Equal(958.22))
			Expect(tareWeightJSON.EmptyWeight).To(Equal(688.89))
			Expect(tareWeightJSON.FullWeight).To(Equal(1700.0))
			Expect(tareWeightJSON.ImageUrl).To(Equal(""))
			Expect(tareWeightJSON.CreatedOn).To(BeTemporally("==", time.Time{}))
			Expect(tareWeightJSON.UpdatedOn).To(BeTemporally(">", time.Time{}))
		})
	})

	Describe("DELETE /v1/tare/weights/:id", func() {
		Describe("with valid input", func() {
			var objectId = bson.NewObjectId()

			BeforeEach(func() {
				body, _ := json.Marshal(gory.Build("tare_weight"))
				request, _ = http.NewRequest("POST", "/v1/tare/weights", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
				collection := session.DB(dbName).C("tare_weights")
				tareWeight := gory.BuildWithParams("tare_weight", gory.Factory{
					"Id": objectId,
				}).(*models.TareWeight)
				collection.Insert(tareWeight)
			})

			It("returns a status code of 200", func() {
				fmt.Printf(recorder.Body.String())
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("does not exist after deleting", func() {
				fmt.Printf(recorder.Body.String())
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))

				recorder = httptest.NewRecorder()
				request, _ = http.NewRequest("GET", "/v1/tare/weights/"+objectId.Hex(), nil)
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))

				var tareWeightJSON models.TareWeight
				json.Unmarshal(recorder.Body.Bytes(), &tareWeightJSON)

				recorder = httptest.NewRecorder()
				request, _ := http.NewRequest("DELETE", "/v1/tare/weights/"+tareWeightJSON.Id.Hex(), nil)
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))

				recorder = httptest.NewRecorder()
				request, _ = http.NewRequest("GET", "/v1/tare/weights/"+tareWeightJSON.Id.Hex(), nil)
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(404))
			})
		})

		Describe("with invalid input", func() {
			BeforeEach(func() {
				json.Marshal(gory.Build("tare_weight"))
				request, _ := http.NewRequest("DELETE", "/v1/tare/weights/"+bson.NewObjectId().Hex(), nil)
				recorder = httptest.NewRecorder()
				server.ServeHTTP(recorder, request)
			})

			It("returns a status code of 404", func() {
				fmt.Printf(recorder.Body.String())
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(404))
			})
		})
	})

	Describe("GET /v1/tare/weights", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/tare/weights?brand=Bombay%20Sapphire", nil)
		})

		Context("when tare weights exist", func() {
			// Insert two valid tare weights into the database
			// before each test in this context.
			BeforeEach(func() {
				collection := session.DB(dbName).C("tare_weights")
				collection.Insert(gory.Build("tare_weight"))
				collection.Insert(gory.Build("tare_weight2"))
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the matching tare weight in the body", func() {
				server.ServeHTTP(recorder, request)

				var tareWeightsJSON []models.TareWeight
				json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
				Expect(len(tareWeightsJSON)).To(Equal(1))

				tareWeightJSON := tareWeightsJSON[0]
				Expect(tareWeightJSON.Brand).To(Equal("Bombay Sapphire"))
				Expect(tareWeightJSON.Category).To(Equal("Liquor"))
				Expect(tareWeightJSON.Name).To(Equal("Bombay Sapphire Gin"))
				Expect(tareWeightJSON.BottleSize).To(Equal(958.22))
				Expect(tareWeightJSON.EmptyWeight).To(Equal(688.89))
				Expect(tareWeightJSON.FullWeight).To(Equal(1700.0))
				Expect(tareWeightJSON.ImageUrl).To(Equal(""))
				Expect(tareWeightJSON.CreatedOn).To(BeTemporally("==", time.Time{}))
				Expect(tareWeightJSON.UpdatedOn).To(BeTemporally("==", time.Time{}))
			})
		})
	})

	Describe("GET /v1/tare/weights", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/tare/weights", nil)
		})

		Context("when no tare weights exist", func() {
			It("returns a status code of 200", func() {
				fmt.Printf(recorder.Body.String())
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns a empty body", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal("[]\n"))
			})
		})

		Context("when tare weights exist", func() {
			// Insert two valid tare weights into the database
			// before each test in this context.
			BeforeEach(func() {
				collection := session.DB(dbName).C("tare_weights")
				collection.Insert(gory.Build("tare_weight"))
				collection.Insert(gory.Build("tare_weight2"))
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those tare weights in the body", func() {
				server.ServeHTTP(recorder, request)

				var tareWeightsJSON []models.TareWeight
				json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
				Expect(len(tareWeightsJSON)).To(Equal(2))

				tareWeightJSON := tareWeightsJSON[0]
				Expect(tareWeightJSON.Brand).To(Equal("Bombay Sapphire"))
				Expect(tareWeightJSON.Category).To(Equal("Liquor"))
				Expect(tareWeightJSON.Name).To(Equal("Bombay Sapphire Gin"))
				Expect(tareWeightJSON.BottleSize).To(Equal(958.22))
				Expect(tareWeightJSON.EmptyWeight).To(Equal(688.89))
				Expect(tareWeightJSON.FullWeight).To(Equal(1700.0))
				Expect(tareWeightJSON.ImageUrl).To(Equal(""))
				Expect(tareWeightJSON.CreatedOn).To(BeTemporally("==", time.Time{}))
				Expect(tareWeightJSON.UpdatedOn).To(BeTemporally("==", time.Time{}))
			})
		})

		Describe("given 25 tare weights", func() {
			BeforeEach(func() {
				collection := session.DB(dbName).C("tare_weights")
				for i := 0; i < 25; i++ {
					if i%2 != 0 {
						collection.Insert(gory.Build("tare_weight"))
					} else {
						collection.Insert(gory.Build("tare_weight2"))
					}
				}
			})

			Context("with no per_page or page query parameters", func() {
				It("returns the first 10 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(10))
				})
			})

			Context("with per_page=15", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "per_page=15"
				})

				It("returns the first 15 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(15))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})

			Context("with invalid per_page=XYZ", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "per_page=XYZ"
				})

				It("returns a status code of 422", func() {
					server.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(422))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})

			Context("with per_page=15&page=2", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "per_page=15&page=2"
				})

				It("returns the second page with 10 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(10))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})

			Context("with page=4 which is out of range", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "page=4"
				})

				It("returns 0 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(0))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})

			Context("with page=-1 which is out of range", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "page=-1"
				})

				It("returns 0 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(0))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})

			Context("with page=X~!$YZ&num_pages=X~!$YZ which is invalid", func() {
				BeforeEach(func() {
					request.URL.RawQuery = "page=X~!$YZ&num_pages=X~!$YZ"
				})

				It("returns 0 tare weights", func() {
					server.ServeHTTP(recorder, request)

					var tareWeightsJSON []models.TareWeight
					json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
					Expect(len(tareWeightsJSON)).To(Equal(0))
				})

				AfterEach(func() {
					request.URL.RawQuery = ""
				})
			})
		})
	})

	Describe("GET /v1/tare/weights/:id", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/tare/weights", nil)
			collection := session.DB(dbName).C("tare_weights")
			collection.Insert(gory.Build("tare_weight"))
			collection.Insert(gory.Build("tare_weight2"))
		})

		Context("when the tare weight exists", func() {
			It("returns a status code of 200", func() {
				fmt.Printf(recorder.Body.String())
				server.ServeHTTP(recorder, request)

				var tareWeightsJSON []models.TareWeight
				json.Unmarshal(recorder.Body.Bytes(), &tareWeightsJSON)
				Expect(len(tareWeightsJSON)).To(Equal(2))

				tareWeightJSON := tareWeightsJSON[0]

				requestOne, _ := http.NewRequest("GET", "/v1/tare/weights/"+tareWeightJSON.Id.Hex(), nil)
				server.ServeHTTP(recorder, requestOne)
				Expect(recorder.Code).To(Equal(200))

				Expect(tareWeightJSON.Brand).To(Equal("Bombay Sapphire"))
				Expect(tareWeightJSON.Category).To(Equal("Liquor"))
				Expect(tareWeightJSON.Name).To(Equal("Bombay Sapphire Gin"))
				Expect(tareWeightJSON.BottleSize).To(Equal(958.22))
				Expect(tareWeightJSON.EmptyWeight).To(Equal(688.89))
				Expect(tareWeightJSON.FullWeight).To(Equal(1700.0))
				Expect(tareWeightJSON.ImageUrl).To(Equal(""))
				Expect(tareWeightJSON.CreatedOn).To(BeTemporally("==", time.Time{}))
				Expect(tareWeightJSON.UpdatedOn).To(BeTemporally("==", time.Time{}))

			})
		})

		Context("when the tare weight does not exist", func() {
			It("returns a empty body", func() {
				request, _ = http.NewRequest("GET", "/v1/tare/weights/507f1f77bcf86cd799439011", nil)
				server.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal("{}\n"))
			})
		})
	})
})
