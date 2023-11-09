package models

import (
	"TodoAPI/initializers"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Name     string `bson:"name,omitempty"`
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
	Todos    []Todo `bson:"todos,omitempty"`
}

func GetUser(email string, user *User) error {
	err := initializers.UsersCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return err
}

func InsertUser(user *User) error {
	_, err := initializers.UsersCollection.InsertOne(context.TODO(), user)
	return err
}

func DeleteUser(email string) error {
	_, err := initializers.UsersCollection.DeleteOne(context.TODO(), bson.M{"email": email})
	return err
}

func UpdateUser(email string, user User) error {
	_, err := initializers.UsersCollection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.D{{"$set", user}})
	return err
}
