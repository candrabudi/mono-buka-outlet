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

	if err := s.seedMidtransSettings(); err != nil {
		return fmt.Errorf("failed to seed midtrans settings: %w", err)
	}

	if err := s.seedAIKnowledgeBase(); err != nil {
		return fmt.Errorf("failed to seed AI knowledge base: %w", err)
	}

	if err := s.seedOutlets(); err != nil {
		return fmt.Errorf("failed to seed outlets: %w", err)
	}

	if err := s.seedEbooks(); err != nil {
		return fmt.Errorf("failed to seed ebooks: %w", err)
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
		{"Affiliator Demo", "affiliator@bukaoutlet.com", "081234567898", "affiliator"},
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

func (s *Seeder) seedMidtransSettings() error {
	log.Println("  🔹 Seeding Midtrans settings...")

	settings := []struct {
		Key         string
		Value       string
		GroupName   string
		Label       string
		Description string
	}{
		{
			Key:         "midtrans_server_key",
			Value:       "SB-Mid-server-XXXXXXXXXXXXXXXXXXXXXXXX",
			GroupName:   "midtrans",
			Label:       "Server Key",
			Description: "Midtrans Server Key (dari dashboard Midtrans)",
		},
		{
			Key:         "midtrans_client_key",
			Value:       "SB-Mid-client-XXXXXXXXXXXXXXXXXXXXXXXX",
			GroupName:   "midtrans",
			Label:       "Client Key",
			Description: "Midtrans Client Key (untuk Snap.js di frontend)",
		},
		{
			Key:         "midtrans_environment",
			Value:       "sandbox",
			GroupName:   "midtrans",
			Label:       "Environment",
			Description: "Environment Midtrans: sandbox atau production",
		},
	}

	for _, st := range settings {
		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM system_settings WHERE key = $1)", st.Key).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("    ⏭️  Setting '%s' already exists, skipping", st.Key)
			continue
		}

		_, err = s.db.Exec(
			`INSERT INTO system_settings (id, key, value, group_name, label, description, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`,
			uuid.New(), st.Key, st.Value, st.GroupName, st.Label, st.Description, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to seed setting %s: %w", st.Key, err)
		}
		log.Printf("    ✅ Setting '%s' seeded", st.Key)
	}

	log.Println("  ✅ Midtrans settings seeded")
	return nil
}
