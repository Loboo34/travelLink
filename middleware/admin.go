package middleware

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Loboo34/travel/auth"
	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	userRepo *auth.UserRepo
}

func NewAdminHandler(userRepo *auth.UserRepo) *AdminService {
	return &AdminService{userRepo: userRepo}
}

func (s *AdminService) Create(ctx context.Context) error {

	var FirstName = os.Getenv("FIRSTNAME")
	var LastName = os.Getenv("LASTNAME")
	var Password = os.Getenv("PASSWORD")
	var Email = os.Getenv("EMAIL")

	existing, err := s.userRepo.FindByEmail(ctx, Email)
	if err != nil {
		return fmt.Errorf("checking email: %w", err)
	}

	if existing != nil {
		return &model.ConflictError{Message: "admin already exixts"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing admin password: %w", err)
	}

	admin := &model.User{
		ID:        primitive.NewObjectID(),
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
		Password:  string(hashedPassword),
		Role:      model.UserRoleAdmin,
		IsActive:  true,
		IsVerified: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, admin); err != nil {
		return fmt.Errorf("creating admin: %w", err)
	}

	return nil

}
