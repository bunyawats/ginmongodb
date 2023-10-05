package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	r := gin.Default()

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	collection := client.Database("test").Collection("people")

	r.POST("/person",
		func(c *gin.Context) {
			//var person Person

			var results map[string]interface{}
			if err := c.ShouldBindJSON(&results); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			insertResult, err := collection.InsertOne(context.TODO(), results)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, gin.H{"message": "Person saved", "id": insertResult.InsertedID})
		})

	r.GET("/person/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		fmt.Println(idStr)

		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		//var person Person
		var results map[string]interface{}
		err = collection.FindOne(
			context.TODO(),
			bson.D{
				{
					Key:   "_id",
					Value: id,
				},
			},
		).Decode(&results)

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		fmt.Println(results)
		c.JSON(http.StatusOK, results)
	})

	r.Run()
}
