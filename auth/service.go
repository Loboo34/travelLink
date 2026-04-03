package auth

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepo
	jwt      *JWTManager
}

func NewUserService(userRepo *UserRepo, jwt *JWTManager) *UserService {
	return &UserService{userRepo: userRepo, jwt: jwt}
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

func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*AuthResult, error) {
	if req.Email == "" {
		return nil, &model.ValidationError{Message: "email is required"}
	}

	if len(req.Password) < 8 {
		return nil, &model.ValidationError{Message: "password must be at least 8 characters"}
	}

	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("checking email: %w", err)
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
		utils.Logger.Error("failed to create user", zap.Error(err))
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
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if user == nil {
		return nil, &model.AuthError{Message: "invalid email or password"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, &model.AuthError{Message: "invalid email or password"}
	}

	if !user.IsActive {
		return nil, &model.AuthError{Message: "account is inactive"}
	}

	token, err := s.jwt.Generate(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, fmt.Errorf("generating tiken: %w", err)
	}

	return &AuthResult{
		Token: token,
		User:  *user,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}

type UpdateProfile struct {
	FirstName   string    `bson:"firstName" json:"firstName"`
	LastName    string    `bson:"lastName" json:"lastName"`
	Gender      string    `bson:"gender" json:"gender"`
	DateOfBirth time.Time `bson:"dateOfBirth" json:"dateOfBirth"`
	Nationality string    `bson:"nationality" json:"nationality"`
	PhoneNumber string    `bson:"phoneNumber" json:"phoneNumber"`
	Email       string    `bson:"email" json:"email"`
}

func (s *UserService) Update(ctx context.Context, userID primitive.ObjectID, req UpdateProfile) (*model.User, error) {

	if err := s.userRepo.UpdateProfile(ctx, userID, req.FirstName, req.LastName, req.Gender, req.Nationality, req.PhoneNumber, req.Email, req.DateOfBirth); err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return &model.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		Nationality: req.Nationality,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}, nil
}
