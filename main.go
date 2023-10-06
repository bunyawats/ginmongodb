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

type (
	//Person struct {
	//	Name string `json:"name"`
	//	Age  int    `json:"age"`
	//}

	PersonMap map[string]interface{}

	ServiceController struct {
		c *mongo.Collection
	}
)

func (sc *ServiceController) getPersonByID(c *gin.Context) {
	idStr := c.Param("id")

	fmt.Println(idStr)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	var personMap PersonMap
	err = sc.c.FindOne(
		context.TODO(),
		bson.D{
			{
				Key:   "_id",
				Value: id,
			},
		},
	).Decode(&personMap)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	fmt.Println(personMap)
	c.JSON(http.StatusOK, personMap)
}

func (sc *ServiceController) creatPerson(c *gin.Context) {

	var personMap PersonMap
	if err := c.ShouldBindJSON(&personMap); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	insertResult, err := sc.c.InsertOne(context.TODO(), personMap)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Person saved", "id": insertResult.InsertedID})
}

func NewServiceController() *ServiceController {

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		panic(err)
	}
	collection := client.Database("test").Collection("people")

	controller := &ServiceController{
		c: collection,
	}

	return controller
}

func main() {

	sc := NewServiceController()

	route := gin.Default()

	route.POST("/person", sc.creatPerson)
	route.GET("/person/:id", sc.getPersonByID)

	_ = route.Run()

}
