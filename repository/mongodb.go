package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	MongoRepository struct {
		*mongo.Client
		*mongo.Collection
	}

	//Person struct {
	//	Name string `json:"name"`
	//	Age  int    `json:"age"`
	//}

	//PersonMap map[string]interface{}
)

func NewMongoRepository(client *mongo.Client) *MongoRepository {

	return &MongoRepository{
		Collection: client.Database("sample_mflix").Collection("movies"),
		Client:     client,
	}
}

func (r MongoRepository) CloseDBConnection() {
	if err := r.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (r MongoRepository) GetMoviesByYear(year int) ([]bson.M, error) {

	filter := bson.D{
		{
			Key:   "year",
			Value: year,
		},
	}

	cursor, err := r.Find(
		context.TODO(),
		filter,
	)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r MongoRepository) GetMoviesById(id primitive.ObjectID) (bson.M, error) {

	filter := bson.D{
		{
			Key:   "_id",
			Value: id,
		},
	}

	var result bson.M
	err := r.FindOne(
		context.TODO(),
		filter,
	).Decode(&result)

	return result, err

}

func (r MongoRepository) CreateNewMovie(payload bson.M) (*mongo.InsertOneResult, error) {

	return r.InsertOne(context.TODO(), payload)
}

func (r MongoRepository) DeleteMovieByID(id primitive.ObjectID) (*mongo.DeleteResult, error) {

	filter := bson.D{
		{
			Key:   "_id",
			Value: id,
		},
	}

	return r.DeleteOne(context.TODO(), filter)
}

func (r MongoRepository) UpdateMovieByID(id primitive.ObjectID, payload bson.M) (*mongo.UpdateResult, error) {

	update := bson.D{
		{"$set", payload},
	}

	return r.UpdateByID(context.TODO(), id, update)

}
