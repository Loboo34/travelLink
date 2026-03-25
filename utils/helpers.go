package utils

import (
	"context"
	"errors"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const (
	contextKeyUserID contextKey = "userID"
	contextKeyRole   contextKey = "role"
)

func GetAdminID(ctx context.Context) (primitive.ObjectID, error) {
	role, err := GetRoleFromContext(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if role != model.UserRoleAdmin {
		return primitive.NilObjectID, errors.New("admin access required")
	}
	return GetUserID(ctx)
}
func GetRoleFromContext(ctx context.Context) (model.UserRole, error) {
	val, ok := ctx.Value(contextKeyRole).(model.UserRole)
	if !ok {
		return "", errors.New("role not found in context")
	}
	return val, nil
}

func GetUserID(ctx context.Context) (primitive.ObjectID, error) {

	val, ok := ctx.Value(contextKeyUserID).(string)
	if !ok || val == "" {
		return primitive.NilObjectID, errors.New("userID not found in context")
	}
	return primitive.ObjectIDFromHex(val)
}
