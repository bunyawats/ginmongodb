package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bunyawats/ginmongodb/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

type (
	GinRoute struct {
		*repository.MongoRepository
		route *gin.Engine
	}
)

func NewGinRoute(r *repository.MongoRepository) *GinRoute {

	ginRoute := &GinRoute{MongoRepository: r}

	route := gin.Default()

	ginRoute.route = route

	route.GET("/movies/year/:year", ginRoute.getMoviesByYear)
	route.GET("/movies/:id", ginRoute.getMovieByID)
	route.POST("/movies", ginRoute.creatNewMovie)
	route.DELETE("/movies/:id", ginRoute.deleteMovieByID)
	route.PUT("/movies/:id", ginRoute.updateMovieById)

	return ginRoute
}

func (sc *GinRoute) getMoviesByYear(c *gin.Context) {

	yearStr := c.Param("year")
	fmt.Printf("request year: %s\n", yearStr)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	results, err := sc.MongoRepository.GetMoviesByYear(year)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fmt.Printf("number of record: %d\n", len(results))
	if len(results) == 0 {
		fmt.Printf("No document was found with the year %d\n", year)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}

	c.JSON(http.StatusOK, results)
}

func (sc *GinRoute) getMovieByID(c *gin.Context) {
	idStr := c.Param("id")

	fmt.Printf("request id: %s\n", idStr)
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	result, err := sc.MongoRepository.GetMoviesById(id)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the id %s\n", idStr)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err == nil {
		fmt.Printf("jsonData = %s\n", jsonData)
	}

	c.JSON(http.StatusOK, result)
}

func (sc *GinRoute) creatNewMovie(c *gin.Context) {

	var reqPayload bson.M
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	insertResult, err := sc.MongoRepository.CreateNewMovie(reqPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, insertResult)
}

func (sc *GinRoute) deleteMovieByID(c *gin.Context) {

	idStr := c.Param("id")
	fmt.Printf("request id: %s\n", idStr)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	deleteResult, err := sc.MongoRepository.DeleteMovieByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, deleteResult)
}

func (sc *GinRoute) updateMovieById(c *gin.Context) {

	idStr := c.Param("id")
	fmt.Printf("request updateMovieById id: %s\n", idStr)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "id binding error"},
		)
		return
	}

	var reqPayload bson.M
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "body binding error"},
		)
		return
	}

	updateResult, err := sc.MongoRepository.UpdateMovieByID(id, reqPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, updateResult)
}

func (sc *GinRoute) Run() interface{} {
	return sc.route.Run()
}
