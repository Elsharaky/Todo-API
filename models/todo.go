package models

import (
	"TodoAPI/initializers"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Todo struct {
	ID          string    `bson:"_id,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Description string    `bson:"description,omitempty"`
	CreatedTime time.Time `bson:"createdTime,omitempty"`
	Status      string    `bson:"status,omitempty"`
}

func CreateTodo(email string, todo Todo) error {
	id := uuid.New()

	todo.CreatedTime = time.Now()
	todo.ID = id.String()

	_, err := initializers.UsersCollection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.M{"$push": bson.M{"todos": todo}})
	return err
}

func GetTodos(email string, status string, todos *[]Todo) error {
	var user User
	err := initializers.UsersCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	for _, todo := range user.Todos {
		if todo.Status == status {
			*todos = append(*todos, todo)
		}
	}

	if status == "" {
		*todos = user.Todos
	}

	return err
}

func DeleteTodo(email, id string) error {
	_, err := initializers.UsersCollection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.M{"$pull": bson.M{"todos": bson.M{"_id": id}}})
	return err
}

func UpdateTodo(email, id string, todo Todo) error {
	updatedFeilds := bson.M{}
	if todo.Title != "" {
		updatedFeilds["todos.$.title"] = todo.Title
	}
	if todo.Description != "" {
		updatedFeilds["todos.$.description"] = todo.Description
	}
	if todo.Status != "" {
		updatedFeilds["todos.$.status"] = todo.Status
	}
	_, err := initializers.UsersCollection.UpdateOne(context.TODO(), bson.M{"email": email, "todos._id": id}, bson.M{"$set": updatedFeilds})
	return err
}
