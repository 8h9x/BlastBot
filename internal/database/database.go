package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database

func GetCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func Init(uri string, database string) error {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db = client.Database(database)

	return nil
}

func Close() error {
	return db.Client().Disconnect(context.Background())
}

func FetchAll[T any](collection string, filter bson.M) ([]T, error) {
	col := GetCollection(collection)

	data := make([]T, 0)

	cursor, err := col.Find(context.Background(), filter)
	if err != nil {
		return data, err
	}

	if err = cursor.All(context.Background(), &data); err != nil {
		return data, err
	}

	return data, nil
}

func Fetch[T any](collection string, filter bson.M) (T, error) {
	col := GetCollection(collection)

	var data T
	if err := col.FindOne(context.Background(), filter).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}
