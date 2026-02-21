package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewAuthUseCase(ur repository.UserRepository, secret string, expiry int) *AuthUseCase {
	return &AuthUseCase{userRepo: ur, jwtSecret: secret, jwtExpiry: expiry}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

func (uc *AuthUseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}
	role := req.Role
	if role == "" {
		role = entity.RoleMitra
	}
	if !entity.IsValidRole(role) {
		return nil, fmt.Errorf("invalid role: %s", role)
	}
	user := &entity.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashed),
		Phone:     req.Phone,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (uc *AuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return uc.userRepo.FindByID(ctx, userID)
}

func (uc *AuthUseCase) GetUsers(ctx context.Context, role string, page, limit int) ([]*entity.User, int, error) {
	return uc.userRepo.FindAll(ctx, role, page, limit)
}

func (uc *AuthUseCase) generateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(uc.jwtExpiry)).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (uc *AuthUseCase) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
		if existing != nil {
			return nil, fmt.Errorf("email sudah terdaftar")
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("gagal hash password")
		}
		user.Password = string(hashed)
	}
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal update user: %w", err)
	}
	return user, nil
}

func (uc *AuthUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	return uc.userRepo.Delete(ctx, id)
}

func (uc *AuthUseCase) ToggleUserActive(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	user.IsActive = !user.IsActive
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal update status: %w", err)
	}
	return user, nil
}
