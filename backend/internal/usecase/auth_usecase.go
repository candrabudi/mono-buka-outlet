package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

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
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Phone           string `json:"phone"`
	Role            string `json:"role"`
}

// validateRegister performs strong validation on registration fields
func validateRegister(req RegisterRequest) []string {
	var errs []string

	// Name validation
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		errs = append(errs, "Nama lengkap minimal 3 karakter")
	}
	if len(name) > 100 {
		errs = append(errs, "Nama lengkap maksimal 100 karakter")
	}

	// Password strength
	if len(req.Password) < 8 {
		errs = append(errs, "Password minimal 8 karakter")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range req.Password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasUpper {
		errs = append(errs, "Password harus mengandung huruf besar")
	}
	if !hasLower {
		errs = append(errs, "Password harus mengandung huruf kecil")
	}
	if !hasDigit {
		errs = append(errs, "Password harus mengandung angka")
	}

	// Confirm password
	if req.Password != req.ConfirmPassword {
		errs = append(errs, "Konfirmasi password tidak cocok")
	}

	// Phone validation (optional but if provided, must be valid)
	if req.Phone != "" {
		phone := strings.ReplaceAll(req.Phone, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phoneRegex := regexp.MustCompile(`^(\+62|62|08)[0-9]{8,13}$`)
		if !phoneRegex.MatchString(phone) {
			errs = append(errs, "Format nomor handphone tidak valid (gunakan format 08xx atau +62xx)")
		}
	}

	return errs
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
	// Strong validation
	if validationErrs := validateRegister(req); len(validationErrs) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(validationErrs, "; "))
	}

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
	if role == entity.RoleAffiliator {
		user.ReferralCode = generateReferralCode()
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

func generateReferralCode() string {
	id := uuid.New()
	code := strings.ToUpper(strings.ReplaceAll(id.String()[:8], "-", ""))
	return "AFF-" + code
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
		if len(strings.TrimSpace(req.Name)) < 3 {
			return nil, fmt.Errorf("nama minimal 3 karakter")
		}
		user.Name = strings.TrimSpace(req.Name)
	}
	if req.Email != "" && req.Email != user.Email {
		existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
		if existing != nil {
			return nil, fmt.Errorf("email sudah terdaftar")
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		phone := strings.ReplaceAll(req.Phone, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phoneRegex := regexp.MustCompile(`^(\+62|62|08)[0-9]{8,13}$`)
		if !phoneRegex.MatchString(phone) {
			return nil, fmt.Errorf("format nomor handphone tidak valid")
		}
		user.Phone = req.Phone
	}
	// Hash and update password if provided
	if req.Password != "" {
		if len(req.Password) < 8 {
			return nil, fmt.Errorf("password minimal 8 karakter")
		}
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

// ChangePasswordRequest is the request body for password change
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

// ChangePassword validates old password and updates to a strong new password
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("password lama salah")
	}

	// Validate new password strength
	if len(req.NewPassword) < 8 {
		return fmt.Errorf("password baru minimal 8 karakter")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range req.NewPassword {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password baru harus mengandung huruf besar, huruf kecil, dan angka")
	}

	// Confirm match
	if req.NewPassword != req.ConfirmNewPassword {
		return fmt.Errorf("konfirmasi password tidak cocok")
	}

	// Check not same as old
	if req.OldPassword == req.NewPassword {
		return fmt.Errorf("password baru tidak boleh sama dengan password lama")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password")
	}
	user.Password = string(hashed)
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("gagal update password: %w", err)
	}
	return nil
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
