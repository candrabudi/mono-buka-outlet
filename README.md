# 🏢 FranchiseHub - Sistem Manajemen Kemitraan Multi-Brand

Sistem manajemen kemitraan/franchise multi-brand yang production-ready, scalable, dan mudah dikembangkan.

## 🏗️ Tech Stack

| Layer     | Technology                  |
| --------- | --------------------------- |
| Backend   | Go 1.22 + Gin Framework     |
| Database  | PostgreSQL 16               |
| Frontend  | Vue 3 + TailwindCSS + Pinia |
| Auth      | JWT (Role-Based Access)     |
| Container | Docker + Docker Compose     |

## 📁 Folder Structure

```
outlet_ready/
├── backend/
│   ├── cmd/api/main.go              # Entry point
│   ├── config/                       # Configuration
│   │   ├── config.go                 # Env loader
│   │   └── database.go              # DB connection
│   ├── internal/
│   │   ├── entity/                   # Domain entities
│   │   ├── repository/              # Repository interfaces
│   │   │   └── postgres/            # PostgreSQL implementations
│   │   ├── usecase/                 # Business logic
│   │   ├── handler/                 # HTTP handlers
│   │   ├── middleware/              # JWT & Role middleware
│   │   ├── router/                  # Route definitions
│   │   ├── migration/              # Migration system
│   │   └── seeder/                  # Database seeder
│   ├── migrations/                   # SQL migration files
│   ├── .env                          # Environment config
│   ├── Dockerfile                    # Docker build
│   └── go.mod                        # Go modules
├── frontend/
│   ├── src/
│   │   ├── layouts/                  # Page layouts
│   │   ├── pages/                   # Page components
│   │   │   ├── public/              # Public pages
│   │   │   ├── auth/                # Auth pages
│   │   │   ├── admin/               # Admin panel
│   │   │   └── mitra/               # Mitra portal
│   │   ├── stores/                  # Pinia stores
│   │   ├── services/                # API services
│   │   ├── router/                  # Vue Router
│   │   └── style.css                # Design system
│   ├── index.html
│   └── vite.config.js
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Option 1: Docker (Recommended)

```bash
docker-compose up -d
```

### Option 2: Manual Setup

#### 1. Database

```bash
createdb franchise_db
```

#### 2. Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials
go mod tidy
go run ./cmd/api -migrate    # Run migrations
go run ./cmd/api -seed       # Seed dummy data
go run ./cmd/api             # Start server on :8080
```

#### 3. Frontend

```bash
cd frontend
npm install
npm run dev                   # Start on :5173
```

### Using Makefile

```bash
make deps          # Install all dependencies
make migrate       # Run migrations
make seed          # Seed database
make dev           # Run backend
make frontend-dev  # Run frontend
make setup         # Fresh migration + seed
```

## 👥 Demo Accounts

| Role          | Email                      | Password    |
| ------------- | -------------------------- | ----------- |
| Super Admin   | admin@franchise.com        | password123 |
| Sales         | sales@franchise.com        | password123 |
| Finance       | finance@franchise.com      | password123 |
| Legal         | legal@franchise.com        | password123 |
| Brand Manager | brandmanager@franchise.com | password123 |
| Mitra         | mitra@franchise.com        | password123 |

## 📊 Database Schema (ERD)

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│    users     │     │   brands    │     │    leads     │
├─────────────┤     ├─────────────┤     ├──────────────┤
│ id (PK)     │     │ id (PK)     │     │ id (PK)      │
│ name        │     │ name        │     │ brand_id (FK)│
│ email       │     │ logo        │     │ sales_id (FK)│
│ password    │     │ description │     │ full_name    │
│ phone       │     │ min_invest  │     │ email        │
│ role        │     │ profit_pct  │     │ phone        │
│ is_active   │     │ est_roi     │     │ status       │
│ created_at  │     │ loc_req     │     │ progress_pct │
│ updated_at  │     │ is_active   │     │ notes        │
│ deleted_at  │     │ created_at  │     │ created_at   │
└──────┬──────┘     └──────┬──────┘     └──────┬───────┘
       │                   │                   │
       │    ┌──────────────┴───────────────────┘
       │    │
┌──────┴────┴──┐     ┌──────────────┐     ┌──────────────┐
│ partnerships │     │   payments   │     │  agreements  │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │     │ id (PK)      │     │ id (PK)      │
│ lead_id (FK) │◄────│ partner_id   │     │ partner_id   │
│ brand_id (FK)│     │ brand_id(FK) │     │ brand_id(FK) │
│ mitra_id(FK) │     │ type         │     │ file_url     │
│ progress_pct │     │ amount       │     │ version      │
│ status       │     │ proof_url    │     │ status       │
│ start_date   │     │ verified_st  │     │ signed_at    │
└──────────────┘     │ verified_by  │     └──────────────┘
       │             └──────────────┘
       │
┌──────┴───────┐     ┌──────────────┐     ┌──────────────┐
│   revenues   │     │activity_logs │     │notifications │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │     │ id (PK)      │     │ id (PK)      │
│ partner_id   │     │ entity_type  │     │ user_id (FK) │
│ brand_id(FK) │     │ entity_id    │     │ title        │
│ month        │     │ action       │     │ message      │
│ revenue      │     │ description  │     │ type         │
│ expense      │     │ old_value    │     │ is_read      │
│ profit       │     │ new_value    │     │ data         │
│ company_share│     │ performed_by │     │ created_at   │
│ mitra_share  │     │ created_at   │     └──────────────┘
└──────────────┘     └──────────────┘

┌──────────────┐
│  locations   │
├──────────────┤
│ id (PK)      │
│ lead_id (FK) │
│ brand_id(FK) │
│ lat, lng     │
│ address      │
│ photo        │
│ approval_st  │
│ survey_notes │
└──────────────┘
```

## 🔐 Role Permission Matrix

| Feature               | Super Admin | Sales | Finance | Legal | Brand Mgr | Mitra |
| --------------------- | :---------: | :---: | :-----: | :---: | :-------: | :---: |
| Dashboard             |     ✅      |  ✅   |   ❌    |  ❌   |    ✅     |  ❌   |
| Brand CRUD            |     ✅      |  ❌   |   ❌    |  ❌   |    ✅     |  ❌   |
| Lead Management       |     ✅      |  ✅   |   ❌    |  ❌   |    ✅     |  ❌   |
| Lead Kanban           |     ✅      |  ✅   |   ❌    |  ❌   |    ✅     |  ❌   |
| Create Partnership    |     ✅      |  ✅   |   ❌    |  ❌   |    ❌     |  ❌   |
| View Partnerships     |     ✅      |  ✅   |   ✅    |  ❌   |    ✅     |  ❌   |
| Create Payment        |     ✅      |  ❌   |   ✅    |  ❌   |    ❌     |  ✅   |
| Verify Payment        |     ✅      |  ❌   |   ✅    |  ❌   |    ❌     |  ❌   |
| Create Agreement      |     ✅      |  ❌   |   ❌    |  ✅   |    ❌     |  ❌   |
| Sign Agreement        |     ✅      |  ❌   |   ❌    |  ✅   |    ❌     |  ✅   |
| Create Revenue        |     ✅      |  ❌   |   ✅    |  ❌   |    ✅     |  ❌   |
| View Own Partnerships |     ❌      |  ❌   |   ❌    |  ❌   |    ❌     |  ✅   |
| View Own Revenue      |     ❌      |  ❌   |   ❌    |  ❌   |    ❌     |  ✅   |

## 📡 API Contract

### Authentication

```
POST /api/v1/auth/login
Body: { "email": "admin@franchise.com", "password": "password123" }
Response: { "success": true, "data": { "token": "jwt...", "user": {...} } }

POST /api/v1/auth/register
Body: { "name": "...", "email": "...", "password": "...", "phone": "...", "role": "mitra" }

GET /api/v1/profile
Headers: Authorization: Bearer <token>
```

### Public

```
GET /api/v1/public/brands?page=1&limit=10
GET /api/v1/public/brands/:id
POST /api/v1/public/consultation
Body: { "brand_id": "uuid", "full_name": "...", "email": "...", "phone": "...", "notes": "..." }
```

### Brands (Admin)

```
GET    /api/v1/brands?page=1&limit=10&active_only=true
POST   /api/v1/brands
GET    /api/v1/brands/:id
PUT    /api/v1/brands/:id
DELETE /api/v1/brands/:id
PATCH  /api/v1/brands/:id/toggle
```

### Leads

```
GET    /api/v1/leads?page=1&limit=20&brand_id=uuid&status=NEW
GET    /api/v1/leads/kanban?brand_id=uuid
POST   /api/v1/leads
GET    /api/v1/leads/:id
PUT    /api/v1/leads/:id
PATCH  /api/v1/leads/:id/status  Body: { "status": "CONSULTATION" }
DELETE /api/v1/leads/:id
```

### Partnerships

```
GET  /api/v1/partnerships?page=1&limit=10&brand_id=uuid&mitra_id=uuid
POST /api/v1/partnerships  Body: { "lead_id": "uuid", "brand_id": "uuid", "mitra_id": "uuid" }
GET  /api/v1/partnerships/:id
GET  /api/v1/partnerships/my  (Mitra only)
```

### Payments

```
POST  /api/v1/payments  Body: { "partnership_id": "uuid", "type": "DP", "amount": 150000000, "proof_url": "..." }
PATCH /api/v1/payments/:id/verify  Body: { "status": "VERIFIED" }
GET   /api/v1/payments/partnership/:partnership_id
```

### Agreements

```
POST  /api/v1/agreements  Body: { "partnership_id": "uuid", "file_url": "..." }
PATCH /api/v1/agreements/:id/sign
GET   /api/v1/agreements/partnership/:partnership_id
```

### Revenue

```
POST /api/v1/revenues  Body: { "partnership_id": "uuid", "month": "2024-01", "revenue": 50000000, "expense": 20000000 }
GET  /api/v1/revenues/partnership/:partnership_id
```

### Dashboard

```
GET /api/v1/dashboard?brand_id=uuid
Response: { "total_leads": 10, "active_mitra": 5, "total_investment": 500000000, "monthly_revenue": 100000000, "leads_by_status": {...}, "revenue_chart": [...] }
```

## 🔄 Lead Status Flow

```
NEW → CONSULTATION → LOCATION_SUBMITTED → SURVEY_APPROVED → MEETING_DONE → READY_FOR_DP → DP_PAID → AGREEMENT_REVIEW → FULLY_PAID → ACTIVE_PARTNERSHIP → RUNNING → COMPLETED
```

## 📈 Partnership Progress (Auto)

| Milestone            | Progress |
| -------------------- | -------- |
| DP Verified          | 25%      |
| Agreement Signed     | 50%      |
| Development/Training | 75%      |
| Running              | 100%     |

## 🎯 Development Milestones

### Phase 1: Foundation ✅

- Clean architecture setup
- Database schema & migrations
- JWT authentication & role middleware
- Seeder with dummy data

### Phase 2: Core Modules ✅

- Brand CRUD + public listing
- Lead management + Kanban board
- Partnership tracking
- Payment with verification + auto progress
- Agreement with signing + auto progress
- Revenue + auto profit sharing

### Phase 3: Frontend ✅

- Vue 3 + TailwindCSS design system
- Public brand listing & detail pages
- Consultation form (auto lead creation)
- Admin dashboard with brand filter
- Lead Kanban board with drag & drop
- Partnership detail with tabs
- Mitra portal with progress tracking

### Phase 4: Production Ready ✅

- Docker + Docker Compose
- Environment configuration
- CORS middleware
- Error handling
- Responsive design

## 📄 License

MIT License
