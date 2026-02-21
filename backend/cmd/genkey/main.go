package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	length := flag.Int("length", 64, "Panjang secret key dalam bytes (default: 64 = 512-bit)")
	format := flag.String("format", "base64", "Format output: base64, hex, atau alphanumeric")
	updateEnv := flag.Bool("update-env", false, "Otomatis update file .env dengan key baru")
	envFile := flag.String("env-file", ".env", "Path ke file .env (default: .env)")
	flag.Parse()

	secret, err := generateSecret(*length, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Gagal generate secret: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔐 JWT Secret Key Generator")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("   Format    : %s\n", *format)
	fmt.Printf("   Length    : %d bytes (%d-bit)\n", *length, *length*8)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("\n   JWT_SECRET=%s\n\n", secret)
	fmt.Println("═══════════════════════════════════════════════════════════════")

	if *updateEnv {
		if err := updateEnvFile(*envFile, secret); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Gagal update .env: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ File %s berhasil diupdate dengan JWT_SECRET baru!\n", *envFile)
	} else {
		fmt.Println("💡 Tip: Tambahkan --update-env untuk otomatis update file .env")
		fmt.Println("   Contoh: go run ./cmd/genkey --update-env")
	}
}

func generateSecret(length int, format string) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("gagal generate random bytes: %w", err)
	}

	switch format {
	case "hex":
		return hex.EncodeToString(bytes), nil
	case "base64":
		return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
	case "alphanumeric":
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		result := make([]byte, length)
		for i, b := range bytes {
			result[i] = charset[int(b)%len(charset)]
		}
		return string(result), nil
	default:
		return "", fmt.Errorf("format tidak valid: %s (gunakan base64, hex, atau alphanumeric)", format)
	}
}

func updateEnvFile(path string, newSecret string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("gagal baca file %s: %w", path, err)
	}

	lines := strings.Split(string(content), "\n")
	found := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "JWT_SECRET=") {
			lines[i] = "JWT_SECRET=" + newSecret
			found = true
			break
		}
	}

	if !found {
		// Cari section JWT dan tambahkan di bawahnya
		for i, line := range lines {
			if strings.TrimSpace(line) == "# JWT" {
				// Sisipkan setelah komentar # JWT
				newLines := make([]string, 0, len(lines)+1)
				newLines = append(newLines, lines[:i+1]...)
				newLines = append(newLines, "JWT_SECRET="+newSecret)
				newLines = append(newLines, lines[i+1:]...)
				lines = newLines
				found = true
				break
			}
		}
	}

	if !found {
		// Jika tidak ada, tambahkan di akhir file
		lines = append(lines, "", "# JWT", "JWT_SECRET="+newSecret)
	}

	output := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(output), 0644)
}
