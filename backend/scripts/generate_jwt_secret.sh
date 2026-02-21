#!/bin/bash
# ═══════════════════════════════════════════════════════════
# 🔐 JWT Secret Key Generator - Bash Version
# ═══════════════════════════════════════════════════════════
# Usage:
#   ./scripts/generate_jwt_secret.sh              # Generate & tampilkan
#   ./scripts/generate_jwt_secret.sh --update      # Generate & update .env
#   ./scripts/generate_jwt_secret.sh --format hex   # Format: base64 (default), hex, alnum
#   ./scripts/generate_jwt_secret.sh --length 64    # Panjang bytes (default: 64)
# ═══════════════════════════════════════════════════════════

set -e

# Default values
LENGTH=64
FORMAT="base64"
UPDATE_ENV=false
ENV_FILE=".env"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --update)    UPDATE_ENV=true; shift ;;
    --format)    FORMAT="$2"; shift 2 ;;
    --length)    LENGTH="$2"; shift 2 ;;
    --env-file)  ENV_FILE="$2"; shift 2 ;;
    --help|-h)
      echo "Usage: $0 [options]"
      echo ""
      echo "Options:"
      echo "  --update           Otomatis update file .env"
      echo "  --format FORMAT    base64 (default), hex, alnum"
      echo "  --length N         Panjang dalam bytes (default: 64 = 512-bit)"
      echo "  --env-file PATH    Path ke .env (default: .env)"
      echo "  --help             Tampilkan bantuan"
      exit 0
      ;;
    *) echo "❌ Unknown option: $1"; exit 1 ;;
  esac
done

# Generate secret berdasarkan format
case "$FORMAT" in
  base64)
    SECRET=$(openssl rand -base64 "$LENGTH" | tr -d '\n' | tr '+/' '-_' | tr -d '=')
    ;;
  hex)
    SECRET=$(openssl rand -hex "$LENGTH")
    ;;
  alnum)
    SECRET=$(openssl rand -base64 "$((LENGTH * 2))" | tr -dc 'a-zA-Z0-9' | head -c "$LENGTH")
    ;;
  *)
    echo "❌ Format tidak valid: $FORMAT (gunakan base64, hex, atau alnum)"
    exit 1
    ;;
esac

BITS=$((LENGTH * 8))

echo ""
echo "🔐 JWT Secret Key Generator"
echo "═══════════════════════════════════════════════════════════════"
echo "   Format    : $FORMAT"
echo "   Length    : ${LENGTH} bytes (${BITS}-bit)"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "   JWT_SECRET=$SECRET"
echo ""
echo "═══════════════════════════════════════════════════════════════"

if [ "$UPDATE_ENV" = true ]; then
  if [ ! -f "$ENV_FILE" ]; then
    echo "❌ File $ENV_FILE tidak ditemukan!"
    exit 1
  fi

  # Backup .env dulu
  cp "$ENV_FILE" "${ENV_FILE}.backup"

  if grep -q "^JWT_SECRET=" "$ENV_FILE"; then
    # Replace existing JWT_SECRET
    sed -i "s|^JWT_SECRET=.*|JWT_SECRET=$SECRET|" "$ENV_FILE"
    echo "✅ File $ENV_FILE berhasil diupdate!"
    echo "📄 Backup disimpan di ${ENV_FILE}.backup"
  else
    # Append JWT_SECRET
    echo "" >> "$ENV_FILE"
    echo "# JWT" >> "$ENV_FILE"
    echo "JWT_SECRET=$SECRET" >> "$ENV_FILE"
    echo "✅ JWT_SECRET ditambahkan ke $ENV_FILE"
  fi
else
  echo "💡 Tip: Tambahkan --update untuk otomatis update file .env"
  echo "   Contoh: $0 --update"
fi

echo ""
