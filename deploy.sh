#!/bin/bash
# ═══════════════════════════════════════════════════════════════
#  BukaOutlet — Production Deploy Script
#  Auto: create DB, migrate, seed, build, deploy
# ═══════════════════════════════════════════════════════════════

set -e

# ┌──────────────────────────────────────────────────────────────
# │ CONFIG — Sesuaikan bagian ini
# └──────────────────────────────────────────────────────────────
REPO_DIR="/home/outlet_ready"
GITHUB_REPO="https://github.com/candrabudi/mono-buka-outlet.git"
BRANCH="master"

# Domain paths
PANEL_DEPLOY="/home/apbo.dinanfarm.com/public_html"
API_DEPLOY="/home/apibo.dinanfarm.com/public_html"
MITRA_DEPLOY="/home/merbo.dinanfarm.com/public_html"

# Domain URLs (untuk .env frontend)
PANEL_DOMAIN="https://apbo.dinanfarm.com"
API_DOMAIN="https://apibo.dinanfarm.com"
MITRA_DOMAIN="https://merbo.dinanfarm.com"

# Backend
PM2_APP_NAME="bukaoutlet-api"
API_PORT="9147"

# Database
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="bukaoutlet_user"
DB_PASS="PasswordSuperKuat_2026!"
DB_NAME="franchise_db"
DB_SSLMODE="disable"

# JWT
JWT_SECRET="MziXzIQePg0Vm84m1J48GpdIyNBEXTxrnRQ09RIQkcfe7GT8E4EksLNi6btUS4ocOtALl5YRpZutCnU2DQq1Tg"

# SMTP
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USERNAME="bagus.candrabudi@gmail.com"
SMTP_PASSWORD="uyfqcvoongkllxft"
SMTP_FROM="bagus.candrabudi@gmail.com"
SMTP_FROM_NAME="BukaOutlet"

# ┌──────────────────────────────────────────────────────────────
# │ HELPERS
# └──────────────────────────────────────────────────────────────
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'
BOLD='\033[1m'

log()  { echo -e "${CYAN}[$(date '+%H:%M:%S')]${NC} $1"; }
ok()   { echo -e "${GREEN}  ✔ $1${NC}"; }
warn() { echo -e "${YELLOW}  ⚠ $1${NC}"; }
fail() { echo -e "${RED}  ✘ $1${NC}"; exit 1; }

echo ""
echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════${NC}"
echo -e "${BOLD}${CYAN}  🚀 BukaOutlet — Production Deploy${NC}"
echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════${NC}"
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 1. GIT CLONE / PULL
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
if [ ! -d "$REPO_DIR/.git" ]; then
  log "📥 Cloning repository..."
  # Hapus dir jika ada tapi bukan git repo (clone gagal sebelumnya)
  if [ -d "$REPO_DIR" ]; then
    rm -rf "$REPO_DIR"
  fi
  git clone -b "$BRANCH" "$GITHUB_REPO" "$REPO_DIR" || fail "Git clone gagal"
  ok "Repository cloned → $REPO_DIR"
else
  log "📥 Pulling latest from $BRANCH..."
  cd "$REPO_DIR"
  git fetch origin "$BRANCH"
  git reset --hard "origin/$BRANCH"
  ok "Code updated to latest $BRANCH"
fi
cd "$REPO_DIR"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 2. DATABASE — Auto create user, db, grant permissions
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "🗄️  Setting up database..."

# Create user if not exists (via postgres superuser)
USER_EXISTS=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" 2>/dev/null || echo "")
if [ "$USER_EXISTS" != "1" ]; then
  sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';" || fail "Gagal membuat DB user"
  ok "User '$DB_USER' created"
else
  ok "User '$DB_USER' already exists — skipped"
fi

# Create database if not exists
DB_EXISTS=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" 2>/dev/null || echo "")
if [ "$DB_EXISTS" != "1" ]; then
  sudo -u postgres psql -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" || fail "Gagal membuat database"
  ok "Database '$DB_NAME' created"
else
  ok "Database '$DB_NAME' already exists — skipped"
fi

# Grant all privileges
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" 2>/dev/null || true
sudo -u postgres psql -d "$DB_NAME" -c "GRANT ALL ON SCHEMA public TO $DB_USER;" 2>/dev/null || true
ok "Permissions granted to '$DB_USER'"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 3. BACKEND .env — Auto generate
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "📝 Generating backend .env..."
mkdir -p "$API_DEPLOY"
cat > "$API_DEPLOY/.env" << ENVEOF
# Application
APP_NAME=franchise-system
APP_ENV=production
APP_PORT=${API_PORT}
APP_URL=${API_DOMAIN}

# Database
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASS}
DB_NAME=${DB_NAME}
DB_SSLMODE=${DB_SSLMODE}

# JWT
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168

# CORS
CORS_ALLOWED_ORIGINS=${PANEL_DOMAIN},${MITRA_DOMAIN}
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# SMTP
SMTP_HOST=${SMTP_HOST}
SMTP_PORT=${SMTP_PORT}
SMTP_USERNAME=${SMTP_USERNAME}
SMTP_PASSWORD=${SMTP_PASSWORD}
SMTP_FROM=${SMTP_FROM}
SMTP_FROM_NAME=${SMTP_FROM_NAME}

# File Upload
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=10485760
ENVEOF
ok "Backend .env generated (port: ${API_PORT})"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 4. BUILD BACKEND
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "🔨 Building backend..."
cd "$REPO_DIR/backend"

# Ensure Go is in PATH (common install locations)
if ! command -v go &>/dev/null; then
  for GO_DIR in /usr/local/go /usr/lib/go /snap/go/current; do
    if [ -x "$GO_DIR/bin/go" ]; then
      export PATH="$GO_DIR/bin:$PATH"
      export GOROOT="$GO_DIR"
      break
    fi
  done
fi
export GOPATH="${GOPATH:-$HOME/go}"
export PATH="$GOPATH/bin:$PATH"

go build -o bukaoutlet-api ./cmd/api || fail "Go build gagal"
ok "Backend binary built ($(go version))"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 5. DEPLOY BACKEND + MIGRATE + SEED
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "📦 Deploying backend..."

# Stop PM2 if running (avoids "Text file busy")
if pm2 describe "$PM2_APP_NAME" > /dev/null 2>&1; then
  pm2 stop "$PM2_APP_NAME" 2>/dev/null || true
  ok "PM2 process stopped"
fi

# Copy binary
cp "$REPO_DIR/backend/bukaoutlet-api" "$API_DEPLOY/bukaoutlet-api"

# Copy migrations
mkdir -p "$API_DEPLOY/migrations"
cp -r "$REPO_DIR/backend/migrations/"* "$API_DEPLOY/migrations/" 2>/dev/null || true

# Copy uploads folder
mkdir -p "$API_DEPLOY/uploads"

ok "Backend files deployed"

# Run migrations (auto-skips already applied)
log "🗄️  Running migrations..."
cd "$API_DEPLOY"
./bukaoutlet-api -migrate 2>&1 | while read line; do echo "  $line"; done || true
ok "Migrations done (new ones applied, existing skipped)"

# Run seed (auto-skips existing data)
log "🌱 Running seeders..."
cd "$API_DEPLOY"
./bukaoutlet-api -seed 2>&1 | while read line; do echo "  $line"; done || true
ok "Seeders done (existing data skipped)"

# Start/Restart with PM2
log "⚡ Starting backend with PM2..."
cd "$API_DEPLOY"

# Kill anything using the port
PORT_PID=$(lsof -ti:$API_PORT 2>/dev/null || true)
if [ -n "$PORT_PID" ]; then
  kill -9 $PORT_PID 2>/dev/null || true
  ok "Killed existing process on port $API_PORT"
  sleep 1
fi

if pm2 describe "$PM2_APP_NAME" > /dev/null 2>&1; then
  pm2 restart "$PM2_APP_NAME"
  ok "Backend restarted (PM2: $PM2_APP_NAME)"
else
  pm2 start ./bukaoutlet-api --name "$PM2_APP_NAME" --cwd "$API_DEPLOY"
  pm2 save
  ok "Backend started (PM2: $PM2_APP_NAME, port: $API_PORT)"
fi

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 5.5 ENSURE NODE.JS 18+ (required by Vite 7)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "🔧 Checking Node.js version..."

# Try to load nvm (common locations)
export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
if [ -s "$NVM_DIR/nvm.sh" ]; then
  . "$NVM_DIR/nvm.sh"
  nvm use 20 2>/dev/null || nvm install 20
  ok "Using Node.js via nvm: $(node -v)"
elif [ -s "/usr/local/nvm/nvm.sh" ]; then
  . "/usr/local/nvm/nvm.sh"
  nvm use 20 2>/dev/null || nvm install 20
  ok "Using Node.js via nvm: $(node -v)"
else
  # No nvm — check current node version
  NODE_MAJOR=$(node -v 2>/dev/null | sed 's/v\([0-9]*\).*/\1/' || echo "0")
  if [ "$NODE_MAJOR" -lt 18 ]; then
    fail "Node.js 18+ required (found: $(node -v 2>/dev/null || echo 'none')). Install nvm and Node 20: curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash && nvm install 20"
  fi
  ok "Node.js $(node -v) — OK (≥18)"
fi

# Also ensure npm is available
command -v npm &>/dev/null || fail "npm not found"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 6. FRONTEND .env — Auto generate
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "📝 Generating frontend .env files..."

# Panel .env
cat > "$REPO_DIR/frontend/panel/.env.production" << ENVEOF
VITE_API_BASE_URL=${API_DOMAIN}/api/v1/admin
ENVEOF
ok "Panel .env.production → ${API_DOMAIN}/api/v1/admin"

# Mitra .env
cat > "$REPO_DIR/frontend/mitra/.env.production" << ENVEOF
VITE_API_BASE_URL=${API_DOMAIN}/api/v1/mitra
ENVEOF
ok "Mitra .env.production → ${API_DOMAIN}/api/v1/mitra"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 7. BUILD & DEPLOY PANEL FRONTEND
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "🎨 Building panel frontend..."
cd "$REPO_DIR/frontend/panel"
npm install --silent 2>/dev/null
PANEL_BUILD_OK=false
if npm run build; then
  ok "Panel built"
  PANEL_BUILD_OK=true
else
  warn "Panel build gagal (Node.js version? Vite error?) — skipping panel deploy"
fi

if [ "$PANEL_BUILD_OK" = true ]; then
  log "📦 Deploying panel..."
  mkdir -p "$PANEL_DEPLOY"
  rm -rf "${PANEL_DEPLOY:?}"/*
  cp -r "$REPO_DIR/frontend/panel/dist/"* "$PANEL_DEPLOY/"

  # .htaccess for Vue Router (OpenLiteSpeed/CyberPanel)
  cat > "$PANEL_DEPLOY/.htaccess" << 'HTEOF'
<IfModule mod_rewrite.c>
  RewriteEngine On
  RewriteBase /
  RewriteRule ^index\.html$ - [L]
  RewriteCond %{REQUEST_FILENAME} !-f
  RewriteCond %{REQUEST_FILENAME} !-d
  RewriteRule . /index.html [L]
</IfModule>
HTEOF
  ok "Panel deployed → $PANEL_DEPLOY (with .htaccess)"
fi

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 8. BUILD & DEPLOY MITRA FRONTEND
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
log "🎨 Building mitra frontend..."
cd "$REPO_DIR/frontend/mitra"
npm install --silent 2>/dev/null
MITRA_BUILD_OK=false
if npm run build; then
  ok "Mitra built"
  MITRA_BUILD_OK=true
else
  warn "Mitra build gagal (Node.js version? Vite error?) — skipping mitra deploy"
fi

if [ "$MITRA_BUILD_OK" = true ]; then
  log "📦 Deploying mitra..."
  mkdir -p "$MITRA_DEPLOY"
  rm -rf "${MITRA_DEPLOY:?}"/*
  cp -r "$REPO_DIR/frontend/mitra/dist/"* "$MITRA_DEPLOY/"

  # .htaccess for Vue Router (OpenLiteSpeed/CyberPanel)
  cat > "$MITRA_DEPLOY/.htaccess" << 'HTEOF'
<IfModule mod_rewrite.c>
  RewriteEngine On
  RewriteBase /
  RewriteRule ^index\.html$ - [L]
  RewriteCond %{REQUEST_FILENAME} !-f
  RewriteCond %{REQUEST_FILENAME} !-d
  RewriteRule . /index.html [L]
</IfModule>
HTEOF
  ok "Mitra deployed → $MITRA_DEPLOY (with .htaccess)"
fi

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# DONE
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo ""
echo -e "${BOLD}${GREEN}══════════════════════════════════════════════════${NC}"
echo -e "${BOLD}${GREEN}  ✅ Deploy Selesai!${NC}"
echo -e "${BOLD}${GREEN}══════════════════════════════════════════════════${NC}"
echo ""
echo -e "  ${BOLD}Panel${NC}   →  ${PANEL_DOMAIN}"
echo -e "  ${BOLD}API${NC}     →  ${API_DOMAIN} (port: ${API_PORT})"
echo -e "  ${BOLD}Mitra${NC}   →  ${MITRA_DOMAIN}"
echo ""
echo -e "  ${BOLD}Database${NC} →  ${DB_NAME}@${DB_HOST}:${DB_PORT}"
echo ""
pm2 status "$PM2_APP_NAME"
echo ""
