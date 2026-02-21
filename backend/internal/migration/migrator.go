package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Init() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (m *Migrator) Up() error {
	if err := m.Init(); err != nil {
		return fmt.Errorf("failed to init migrations table: %w", err)
	}

	files, err := m.getMigrationFiles("up")
	if err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, file := range files {
		version := m.extractVersion(file)
		if _, ok := applied[version]; ok {
			continue
		}

		log.Printf("🔄 Applying migration: %s", file)
		content, err := os.ReadFile(filepath.Join(m.migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		log.Printf("✅ Applied migration: %s", file)
	}

	return nil
}

func (m *Migrator) Down() error {
	if err := m.Init(); err != nil {
		return fmt.Errorf("failed to init migrations table: %w", err)
	}

	files, err := m.getMigrationFiles("down")
	if err != nil {
		return err
	}

	// Reverse order for down migrations
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, file := range files {
		version := m.extractVersion(file)
		if _, ok := applied[version]; !ok {
			continue
		}

		log.Printf("🔄 Reverting migration: %s", file)
		content, err := os.ReadFile(filepath.Join(m.migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to revert migration %s: %w", file, err)
		}

		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration revert %s: %w", file, err)
		}

		log.Printf("✅ Reverted migration: %s", file)
	}

	return nil
}

func (m *Migrator) Fresh() error {
	log.Println("🔄 Running fresh migration (drop all + migrate)")
	if err := m.Down(); err != nil {
		// Ignore errors on fresh
		log.Printf("⚠️ Down migration had issues (this is normal for fresh): %v", err)
	}

	// Drop migration tracker too
	m.db.Exec("DROP TABLE IF EXISTS schema_migrations")

	return m.Up()
}

func (m *Migrator) getMigrationFiles(direction string) ([]string, error) {
	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	suffix := fmt.Sprintf(".%s.sql", direction)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return map[string]bool{}, nil
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

func (m *Migrator) extractVersion(filename string) string {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}
