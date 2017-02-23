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

	Describe("PUT /v1/tare/weights/:id", func() {
		BeforeEach(func() {
			collection := session.DB(dbName).C("tare_weights")

			objectId := bson.NewObjectId()

			tareWeight2 := gory.BuildWithParams("tare_weight2", gory.Factory{
				"Id":        objectId,
				"CreatedOn": time.Now().UnixNano() / int64(time.Millisecond),
				"UpdatedOn": time.Now().UnixNano() / int64(time.Millisecond),
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

			// @todo Get the specific record instead
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest("GET", "/v1/tare/weights", nil)
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))

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
			Expect(tareWeightJSON.CreatedOn).To(Equal(int64(0)))
			Expect(tareWeightJSON.UpdatedOn).NotTo(Equal(int64(0)))
		})
	})

	Describe("DELETE /v1/tare/weights/:id", func() {
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

	Describe("POST /v1/tare/weights", func() {
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
				Expect(tareWeightJSON.CreatedOn).To(Equal(int64(0)))
				Expect(tareWeightJSON.UpdatedOn).To(Equal(int64(0)))
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
			// Insert two valid signatures into the database
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
				Expect(tareWeightJSON.CreatedOn).To(Equal(int64(0)))
				Expect(tareWeightJSON.UpdatedOn).To(Equal(int64(0)))
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
				Expect(tareWeightJSON.CreatedOn).To(Equal(int64(0)))
				Expect(tareWeightJSON.UpdatedOn).To(Equal(int64(0)))

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
