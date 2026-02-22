#!/bin/bash
# ═══════════════════════════════════════════════════
#  BukaOutlet — Deploy Script
#  Pulls latest code, builds, deploys to production
# ═══════════════════════════════════════════════════

set -e

# ── Config ──────────────────────────────────────────
REPO_DIR="/home/outlet_ready"                    # clone directory di server
PANEL_DEPLOY="/home/apbo.dinanfarm.com"          # panel frontend
API_DEPLOY="/home/apibo.dinanfarm.com"           # backend API
MITRA_DEPLOY="/home/merbo.dinanfarm.com"         # mitra frontend
PM2_APP_NAME="bukaoutlet-api"                    # nama proses PM2
API_PORT="8080"                                  # port backend

# ── Colors ──────────────────────────────────────────
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

log()  { echo -e "${CYAN}[$(date '+%H:%M:%S')]${NC} $1"; }
ok()   { echo -e "${GREEN}  ✔ $1${NC}"; }
warn() { echo -e "${YELLOW}  ⚠ $1${NC}"; }
fail() { echo -e "${RED}  ✘ $1${NC}"; exit 1; }

echo ""
echo -e "${CYAN}══════════════════════════════════════════════${NC}"
echo -e "${CYAN}  🚀 BukaOutlet Deploy${NC}"
echo -e "${CYAN}══════════════════════════════════════════════${NC}"
echo ""

# ── 1. Pull Latest Code ────────────────────────────
log "📥  Pulling latest code..."
cd "$REPO_DIR"
git pull origin master || fail "Git pull gagal"
ok "Code updated"

# ── 2. Build Backend (Go) ──────────────────────────
log "🔨  Building backend..."
cd "$REPO_DIR/backend"
go build -o bukaoutlet-api ./cmd/api || fail "Go build gagal"
ok "Backend binary built"

# ── 3. Build Panel Frontend ────────────────────────
log "🎨  Building panel frontend..."
cd "$REPO_DIR/frontend/panel"
npm install --silent 2>/dev/null
npm run build || fail "Panel build gagal"
ok "Panel frontend built"

# ── 4. Build Mitra Frontend ────────────────────────
log "🎨  Building mitra frontend..."
cd "$REPO_DIR/frontend/mitra"
npm install --silent 2>/dev/null
npm run build || fail "Mitra build gagal"
ok "Mitra frontend built"

# ── 5. Deploy Backend API ──────────────────────────
log "📦  Deploying backend API..."

# Copy binary
cp "$REPO_DIR/backend/bukaoutlet-api" "$API_DEPLOY/bukaoutlet-api"

# Copy .env if not exists
if [ ! -f "$API_DEPLOY/.env" ]; then
  cp "$REPO_DIR/backend/.env" "$API_DEPLOY/.env"
  warn ".env copied — pastikan config production sudah benar!"
fi

# Copy migrations
mkdir -p "$API_DEPLOY/migrations"
cp -r "$REPO_DIR/backend/migrations/"* "$API_DEPLOY/migrations/" 2>/dev/null || true

# Copy uploads folder structure
mkdir -p "$API_DEPLOY/uploads"

# Restart with PM2
cd "$API_DEPLOY"
if pm2 describe "$PM2_APP_NAME" > /dev/null 2>&1; then
  pm2 restart "$PM2_APP_NAME"
  ok "Backend restarted via PM2"
else
  pm2 start ./bukaoutlet-api --name "$PM2_APP_NAME" --cwd "$API_DEPLOY"
  pm2 save
  ok "Backend started via PM2 (new process)"
fi

# ── 6. Deploy Panel Frontend ──────────────────────
log "📦  Deploying panel frontend..."
mkdir -p "$PANEL_DEPLOY"
rm -rf "$PANEL_DEPLOY"/*
cp -r "$REPO_DIR/frontend/panel/dist/"* "$PANEL_DEPLOY/"
ok "Panel deployed to $PANEL_DEPLOY"

# ── 7. Deploy Mitra Frontend ──────────────────────
log "📦  Deploying mitra frontend..."
mkdir -p "$MITRA_DEPLOY"
rm -rf "$MITRA_DEPLOY"/*
cp -r "$REPO_DIR/frontend/mitra/dist/"* "$MITRA_DEPLOY/"
ok "Mitra deployed to $MITRA_DEPLOY"

# ── 8. Run Migrations ─────────────────────────────
log "🗄️  Running database migrations..."
cd "$API_DEPLOY"
./bukaoutlet-api -migrate &
MIGRATE_PID=$!
sleep 3
kill $MIGRATE_PID 2>/dev/null || true
ok "Migrations executed"

# ── Done! ─────────────────────────────────────────
echo ""
echo -e "${GREEN}══════════════════════════════════════════════${NC}"
echo -e "${GREEN}  ✅ Deploy selesai!${NC}"
echo -e "${GREEN}══════════════════════════════════════════════${NC}"
echo ""
echo -e "  Panel  →  https://apbo.dinanfarm.com"
echo -e "  API    →  https://apibo.dinanfarm.com"
echo -e "  Mitra  →  https://merbo.dinanfarm.com"
echo ""
echo -e "  PM2 status:"
pm2 status "$PM2_APP_NAME"
echo ""
