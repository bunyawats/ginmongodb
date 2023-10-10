package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r MongoRepository) CloseDBConnection() {
	if err := r.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (r MongoRepository) GetAllMovies(year int) ([]bson.M, error) {

	filter := bson.D{
		{"year", year},
	}

	cursor, err := r.c.Find(
		context.TODO(),
		filter,
	)

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results, err
}

func (r MongoRepository) GetMoviesById(id primitive.ObjectID) (bson.M, error) {

	var result bson.M

	err := r.c.FindOne(
		context.TODO(),
		bson.D{
			{
				Key:   "_id",
				Value: id,
			},
		},
	).Decode(&result)

	return result, err

}

func (r MongoRepository) CreateNewMovie(payload bson.M) (*mongo.InsertOneResult, error) {

	return r.c.InsertOne(context.TODO(), payload)
}
