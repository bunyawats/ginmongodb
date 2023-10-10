package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	MongoRepository struct {
		client *mongo.Client
		c      *mongo.Collection
	}

	//Person struct {
	//	Name string `json:"name"`
	//	Age  int    `json:"age"`
	//}

	//PersonMap map[string]interface{}
)

func NewMongoRepository(client *mongo.Client) *MongoRepository {

	collection := client.Database("sample_mflix").Collection("movies")

	return &MongoRepository{
		c:      collection,
		client: client,
	}
}

func (r MongoRepository) GetAllMovies(title string) (bson.M, error) {

	var result bson.M

	err := r.c.FindOne(
		context.TODO(),
		bson.D{
			{"title", title},
		},
	).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the title %s\n", title)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r MongoRepository) CloseDBConnection() {
	if err := r.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

//func (sc *ServiceController) getPersonByID(c *gin.Context) {
//	idStr := c.Param("id")
//
//	fmt.Println(idStr)
//
//	id, err := primitive.ObjectIDFromHex(idStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest,
//			gin.H{"error": err.Error()},
//		)
//		return
//	}
//
//	var personMap PersonMap
//	err = sc.c.FindOne(
//		context.TODO(),
//		bson.D{
//			{
//				Key:   "_id",
//				Value: id,
//			},
//		},
//	).Decode(&personMap)
//
//	if err != nil {
//		c.AbortWithStatus(http.StatusNotFound)
//		return
//	}
//
//	fmt.Println(personMap)
//	c.JSON(http.StatusOK, personMap)
//}
//
//func (sc *ServiceController) creatPerson(c *gin.Context) {
//
//	var personMap PersonMap
//	if err := c.ShouldBindJSON(&personMap); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	insertResult, err := sc.c.InsertOne(context.TODO(), personMap)
//	if err != nil {
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(200, gin.H{"message": "Person saved", "id": insertResult.InsertedID})
//}
