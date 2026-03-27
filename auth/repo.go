package auth

import (
	"context"
	"errors"
	"fmt"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, register *model.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, register)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("database error finding user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, userID primitive.ObjectID) error{
	_, err := r.db.Collection("users").DeleteOne(ctx, bson.M{"_id": userID })
	if err != nil{
		return fmt.Errorf("deleting user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error){
	var user model.User

	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil, fmt.Errorf("usr not found: %w", err)
		}
		return nil, fmt.Errorf("database err: %w", err)
	}

	return &user, nil
}
