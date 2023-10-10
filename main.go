package main

import (
	"context"
	"github.com/bunyawats/ginmongodb/repository"
	"github.com/bunyawats/ginmongodb/restapi"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
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

	ginRoute := restapi.NewGinRoute(repo)

	_ = ginRoute.Run()

}
