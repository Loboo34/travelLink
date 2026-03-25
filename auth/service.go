package auth

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepo
	jwt      *JWTManager
}

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *UserService) Resgister(ctx context.Context, req RegisterRequest) (*AuthResult, error) {
	if req.Email == "" {
		return nil, &model.ValidationError{Message: "email is required"}
	}

	if len(req.Password) < 8 {
		return nil, &model.ValidationError{Message: "password must be at least 8 charachters"}
	}

	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error fnding email %w", err)
	}

	if existing != nil {
		return nil, &model.ConflictError{Message: "email already registered"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password")
	}

	userID := primitive.NewObjectID()
	now := time.Now()

	user := &model.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hash),
		Role:      model.UserRoleUser,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	token, err := s.jwt.Generate(userID, user.Role, req.Email)
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	return &AuthResult{
		Token: token,
		User:  *user}, nil

}

func (s *UserService) Login(ctx context.Context, req LoginRequest) (*AuthResult, error) {
	if req.Email == "" {
		return nil, &model.ValidationError{Message: "email is required"}
	}

	if req.Password == "" {
		return nil, &model.ValidationError{Message: "password is required"}
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil{
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if user == nil{
		return nil, &model.AuthError{Message: "invalid email or password"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil{
		return nil, &model.AuthError{Message: "invalid email or password"}
	}

	if !user.IsActive{
		return nil, &model.AuthError{Message: "account is inactive"}
	}

	token, err := s.jwt.Generate(user.ID, user.Role, user.Email)
	if err != nil{
		return nil, fmt.Errorf("generating tiken: %w", err )
	}

	return &AuthResult{
		Token: token,
		User: *user,
	}, nil
}
