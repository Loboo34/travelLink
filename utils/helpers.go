package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetAdminID() (primitive.ObjectID, error) {

	id := primitive.NewObjectID()
	return id, nil
}

func GetUserID() (primitive.ObjectID, error) {

	id := primitive.NewObjectID()
	return id, nil
}