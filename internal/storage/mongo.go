package storage

import (
	"context"
	"fmt"
	"toDoList/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStorage struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoConnection(uri, dbName, collectionName string) (*mongoStorage, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongo: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not ping mongo: %v", err)
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &mongoStorage{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (m *mongoStorage) AddTodo(todo model.ToDo) error {
	if todo.ID == "" {
		todo.ID = primitive.NewObjectID().Hex()
	}
	_, err := m.collection.InsertOne(context.Background(), todo)
	return err
}

func (m *mongoStorage) GetTodos() ([]model.ToDo, error) {
	cursor, err := m.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var todos []model.ToDo
	for cursor.Next(context.Background()) {
		var todo model.ToDo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (m *mongoStorage) GetTodoById(id string) (model.ToDo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.ToDo{}, err
	}
	var todo model.ToDo
	err = m.collection.FindOne(context.Background(), bson.D{{"_id", objectID}}).Decode(&todo)
	if err != nil {
		return model.ToDo{}, err
	}
	return todo, nil
}

func (m *mongoStorage) UpdateTodo(id string, todo model.ToDo) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.collection.UpdateOne(
		context.Background(),
		bson.D{{"_id", objectID}},
		bson.D{{"$set", bson.D{{"title", todo.Title}, {"status", todo.Status}}}},
	)
	return err
}

func (m *mongoStorage) DeleteTodo(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(context.Background(), bson.D{{"_id", objectID}})
	return err
}

func (m *mongoStorage) Close() {
	m.client.Disconnect(context.Background())
}
