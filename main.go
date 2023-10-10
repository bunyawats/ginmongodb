package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bunyawats/ginmongodb/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
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

	//route.POST("/person", sc.creatPerson)
	//route.GET("/person/:id", sc.getPersonByID)
	route.GET("/movies", sc.getAllMovie)

	_ = route.Run()

}

func (sc *ServiceController) getAllMovie(c *gin.Context) {

	title := "Back to the Future"

	result, err := sc.MongoRepository.GetAllMovies(title)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the title %s\n", title)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fmt.Printf("result = %s\n", result)
	c.JSON(http.StatusOK, result)
}
