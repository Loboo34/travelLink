package auth

import (
	"context"
	"errors"
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

func (r *RegisterRequest) Validate() error {
	if r.FirstName == "" {
		return errors.New("first name required")
	}
	if r.LastName == "" {
		return errors.New("first name required")
	}
	if len(r.Password) < 8 {
		return errors.New("first name required")
	}
	if r.Email == "" {
		return errors.New("first name required")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Password == "" {
		return errors.New("first name required")
	}
	if r.Email == "" {
		return errors.New("first name required")
	}
	return nil
}

type AuthResult struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *UserService) Register(ctx context.Context, req RegisterRequest) (*AuthResult, error) {

	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
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

	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
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

func (r *UpdateProfile) Validate() error {
	if r.FirstName == "" {
		return errors.New("first name required")
	}
	if r.LastName == "" {
		return errors.New("first name required")
	}
	if r.Gender == "" {
		return errors.New("gender required")
	}
	if r.Nationality == "" {
		return errors.New("nationality required")
	}

	if len(r.PhoneNumber) < 10 {
		return errors.New("phone number is less than 10 numbers")
	}

	return nil
}

func (s *UserService) Update(ctx context.Context, userID primitive.ObjectID, req UpdateProfile) (*model.User, error) {

	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
	}

	if err := s.userRepo.UpdateProfile(ctx, userID, req.FirstName, req.LastName, req.Gender, req.Nationality, req.PhoneNumber, req.Email, req.DateOfBirth); err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	updated, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching updated user: %w", err)
	}
	return updated, nil
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (s *UserService) ChangePassword(ctx context.Context, userID primitive.ObjectID, req ChangePasswordRequest) error {
	if req.CurrentPassword == "" {
		return &model.ValidationError{Message: "current password is required"}
	}
	if len(req.NewPassword) < 8 {
		return &model.ValidationError{Message: "new password must be at least 8 characters"}
	}
	if req.CurrentPassword == req.NewPassword {
		return &model.ValidationError{Message: "new password must differ from current password"}
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return &model.AuthError{Message: "current password is incorrect"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

func (s *UserService) ProfilePic(ctx context.Context, userID primitive.ObjectID, image string) error {

	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.userRepo.ProfilePic(ctx, userID, image); err != nil {
		return fmt.Errorf("updating profile picture")
	}

	return nil
}
