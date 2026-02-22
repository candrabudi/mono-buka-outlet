package seeder

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run() error {
	log.Println("🌱 Running seeders...")

	if err := s.seedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	log.Println("✅ Seeding completed successfully")
	return nil
}

func (s *Seeder) seedUsers() error {
	log.Println("  🔹 Seeding users...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []struct {
		Name  string
		Email string
		Phone string
		Role  string
	}{
		{"Master Admin", "bagus.candrabudi@gmail.com", "081234567890", "master"},
		{"Admin Outlet", "admin@bukaoutlet.com", "081234567891", "admin"},
		{"Finance Officer", "finance@bukaoutlet.com", "081234567892", "finance"},
		{"Mitra Demo", "mitra@bukaoutlet.com", "081234567895", "mitra"},
		{"Ahmad Rizki", "ahmad@bukaoutlet.com", "081234567896", "admin"},
		{"Siti Nurhaliza", "siti@bukaoutlet.com", "081234567897", "mitra"},
	}

	for _, u := range users {
		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", u.Email).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		_, err = s.db.Exec(
			`INSERT INTO users (id, name, email, password, phone, role, is_active, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, true, $7, $7)`,
			uuid.New(), u.Name, u.Email, string(hashedPassword), u.Phone, u.Role, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", u.Email, err)
		}
	}

	log.Println("  ✅ Users seeded")
	return nil
}
