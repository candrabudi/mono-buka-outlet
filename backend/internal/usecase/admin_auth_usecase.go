package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/franchise-system/backend/internal/repository/postgres"
	"github.com/franchise-system/backend/internal/service/email"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminAuthUseCase struct {
	userRepo     repository.UserRepository
	otpRepo      *postgres.OTPRepo
	emailService *email.EmailService
	jwtSecret    string
	jwtExpiry    int
}

func NewAdminAuthUseCase(
	ur repository.UserRepository,
	otpRepo *postgres.OTPRepo,
	emailSvc *email.EmailService,
	secret string,
	expiry int,
) *AdminAuthUseCase {
	return &AdminAuthUseCase{
		userRepo:     ur,
		otpRepo:      otpRepo,
		emailService: emailSvc,
		jwtSecret:    secret,
		jwtExpiry:    expiry,
	}
}

// ── Request / Response types ──

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type AdminAuthResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ── Step 1: Login (validate credentials, send OTP) ──

func (uc *AdminAuthUseCase) Login(ctx context.Context, req AdminLoginRequest) (*AdminLoginResponse, error) {
	// Find user
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("email atau password salah")
	}

	// Check active
	if !user.IsActive {
		return nil, fmt.Errorf("akun tidak aktif")
	}

	// Check admin role
	if !entity.IsAdminRole(user.Role) {
		return nil, fmt.Errorf("akun tidak memiliki akses admin")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("email atau password salah")
	}

	// Generate OTP
	code := generateOTPCode()

	// Invalidate previous OTPs
	_ = uc.otpRepo.InvalidateAll(ctx, user.Email, entity.OTPPurposeAdminLogin)

	// Create OTP record
	otp := &entity.OTPCode{
		ID:        uuid.New(),
		UserID:    user.ID,
		Email:     user.Email,
		Code:      code,
		Purpose:   entity.OTPPurposeAdminLogin,
		ExpiresAt: time.Now().Add(time.Duration(entity.OTPExpiryMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("gagal membuat OTP: %w", err)
	}

	// Send OTP via email
	if err := uc.emailService.SendOTP(user.Email, code, entity.OTPPurposeAdminLogin); err != nil {
		return nil, fmt.Errorf("gagal mengirim OTP ke email: %w", err)
	}

	return &AdminLoginResponse{
		SessionID: otp.ID.String(),
		Email:     user.Email,
		Message:   fmt.Sprintf("Kode OTP telah dikirim ke %s", maskEmail(user.Email)),
	}, nil
}

// ── Step 2: Verify OTP → get JWT ──

func (uc *AdminAuthUseCase) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*AdminAuthResponse, error) {
	// Find OTP
	otp, err := uc.otpRepo.FindByCode(ctx, req.Email, req.Code, entity.OTPPurposeAdminLogin)
	if err != nil {
		return nil, fmt.Errorf("kode OTP tidak valid")
	}

	// Check expiry
	if otp.IsExpired() {
		return nil, fmt.Errorf("kode OTP sudah kedaluwarsa")
	}

	// Check already used
	if otp.IsUsed() {
		return nil, fmt.Errorf("kode OTP sudah digunakan")
	}

	// Mark as used
	if err := uc.otpRepo.MarkUsed(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("gagal memverifikasi OTP")
	}

	// Find user
	user, err := uc.userRepo.FindByID(ctx, otp.UserID)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}

	// Generate JWT
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token")
	}

	return &AdminAuthResponse{Token: token, User: user}, nil
}

// ── Resend OTP ──

func (uc *AdminAuthUseCase) ResendOTP(ctx context.Context, req ResendOTPRequest) (*AdminLoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("email tidak terdaftar")
	}
	if !entity.IsAdminRole(user.Role) {
		return nil, fmt.Errorf("akun tidak memiliki akses admin")
	}

	// Invalidate previous
	_ = uc.otpRepo.InvalidateAll(ctx, user.Email, entity.OTPPurposeAdminLogin)

	// Generate new OTP
	code := generateOTPCode()
	otp := &entity.OTPCode{
		ID:        uuid.New(),
		UserID:    user.ID,
		Email:     user.Email,
		Code:      code,
		Purpose:   entity.OTPPurposeAdminLogin,
		ExpiresAt: time.Now().Add(time.Duration(entity.OTPExpiryMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("gagal membuat OTP: %w", err)
	}
	if err := uc.emailService.SendOTP(user.Email, code, entity.OTPPurposeAdminLogin); err != nil {
		return nil, fmt.Errorf("gagal mengirim OTP: %w", err)
	}

	return &AdminLoginResponse{
		SessionID: otp.ID.String(),
		Email:     user.Email,
		Message:   fmt.Sprintf("Kode OTP baru telah dikirim ke %s", maskEmail(user.Email)),
	}, nil
}

// ── Profile ──

func (uc *AdminAuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return uc.userRepo.FindByID(ctx, userID)
}

// ── Helpers ──

func (uc *AdminAuthUseCase) generateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"panel":   "admin",
		"exp":     time.Now().Add(time.Hour * time.Duration(uc.jwtExpiry)).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func generateOTPCode() string {
	code := ""
	for i := 0; i < entity.OTPLength; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code
}

func maskEmail(email string) string {
	at := 0
	for i, c := range email {
		if c == '@' {
			at = i
			break
		}
	}
	if at <= 2 {
		return "***" + email[at:]
	}
	return email[:2] + "***" + email[at:]
}
