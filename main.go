package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bunyawats/ginmongodb/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strconv"
)

type (
	ServiceController struct {
		*repository.MongoRepository
	}
)

var (
	repo   *repository.MongoRepository
	client *mongo.Client
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		panic(err)
	}

	repo = repository.NewMongoRepository(client)

}

func main() {

	defer func() {
		repo.CloseDBConnection()
	}()

	route := gin.Default()

	sc := &ServiceController{
		repo,
	}

	route.GET("/movies/year/:year", sc.getAllMovie)
	route.GET("/movies/:id", sc.getMovieByID)
	route.POST("/movies", sc.creatNewMovie)
	route.DELETE("/movies/:id", sc.deleteMovieByID)

	_ = route.Run()

}

func (sc *ServiceController) getAllMovie(c *gin.Context) {

	yearStr := c.Param("year")
	fmt.Printf("request year: %s\n", yearStr)

	fmt.Println(yearStr)
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	results, err := sc.MongoRepository.GetAllMovies(year)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the year %s\n", year)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
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

func (sc *ServiceController) getMovieByID(c *gin.Context) {
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

func (sc *ServiceController) creatNewMovie(c *gin.Context) {

	var reqPayload bson.M
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	insertResult, err := sc.MongoRepository.CreateNewMovie(reqPayload)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insertResult)
}

func (sc *ServiceController) deleteMovieByID(c *gin.Context) {

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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deleteResult)
}
