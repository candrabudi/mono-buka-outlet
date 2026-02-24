# 🤖 AI Konsultan BukaOutlet — System Flow

## Arsitektur Lengkap

```
┌─────────────────────────────────────────────────────────────────────┐
│                        FRONTEND (Vue.js)                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────────┐    ┌──────────────────────────────┐  │
│  │   Portal Mitra            │    │   Admin Panel                 │  │
│  │   (ChatWidget.vue)        │    │   (AISettings.vue)            │  │
│  │                           │    │                               │  │
│  │  • Floating chat button   │    │  • Kelola Knowledge Base      │  │
│  │  • Kirim pesan            │    │  • Edit System Prompt         │  │
│  │  • Riwayat percakapan     │    │  • Konfigurasi Model AI       │  │
│  │  • Quick actions          │    │  • Monitor penggunaan         │  │
│  │  • Markdown rendering     │    │                               │  │
│  └──────────┬───────────────┘    └────────────┬──────────────────┘  │
│             │                                  │                    │
└─────────────┼──────────────────────────────────┼────────────────────┘
              │ POST /api/v1/mitra/chat          │ CRUD /api/v1/admin/ai/*
              │ GET  /api/v1/mitra/chat/...      │
              ▼                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        BACKEND (Go + Gin)                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                      CHAT HANDLER                             │   │
│  │  POST /chat          → Kirim pesan & terima respons           │   │
│  │  GET  /conversations → Daftar percakapan user                 │   │
│  │  GET  /conversations/:id → Pesan dalam percakapan             │   │
│  │  DELETE /conversations/:id → Hapus percakapan                 │   │
│  └──────────────────┬───────────────────────────────────────────┘   │
│                     │                                               │
│                     ▼                                               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                   CHATBOT SERVICE                             │   │
│  │                                                               │   │
│  │  1. resolveConversation() → Buat/ambil conversation           │   │
│  │  2. saveUserMessage()     → Simpan pesan user ke DB           │   │
│  │  3. buildGPTMessages()    → Rakit konteks untuk GPT           │   │
│  │  4. openai.ChatCompletion() → Kirim ke OpenAI API             │   │
│  │  5. saveAssistantMessage() → Simpan respons ke DB             │   │
│  │  6. generateQuickActions() → Saran tindak lanjut              │   │
│  │                                                               │   │
│  │  Cache Layer (5 menit):                                       │   │
│  │  • cachedPrompt   → System prompt dari DB                     │   │
│  │  • cachedKB       → Knowledge base dari DB                    │   │
│  │  • cachedLiveData → Data outlet/paket/ebook dari DB           │   │
│  └──────────────────┬───────────────────────────────────────────┘   │
│                     │                                               │
│                     ▼                                               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              buildGPTMessages() Detail                        │   │
│  │                                                               │   │
│  │  Messages Array yang dikirim ke OpenAI:                       │   │
│  │                                                               │   │
│  │  ┌────────────────────────────────────────────┐              │   │
│  │  │ 1. SYSTEM PROMPT (dari ai_system_prompts)   │              │   │
│  │  │    • Identitas AI Konsultan                 │              │   │
│  │  │    • Aturan ketat (hanya bisnis outlet)     │              │   │
│  │  │    • Format jawaban                         │              │   │
│  │  │    • Gaya komunikasi                        │              │   │
│  │  └────────────────────────────────────────────┘              │   │
│  │  ┌────────────────────────────────────────────┐              │   │
│  │  │ 2. KNOWLEDGE BASE (dari ai_knowledge_base)  │              │   │
│  │  │    • 25+ artikel referensi                  │              │   │
│  │  │    • 11 kategori topik                      │              │   │
│  │  │    • Dikelompokkan per kategori             │              │   │
│  │  └────────────────────────────────────────────┘              │   │
│  │  ┌────────────────────────────────────────────┐              │   │
│  │  │ 3. LIVE BUSINESS DATA (dari tabel DB)       │              │   │
│  │  │    • outlets → Nama, harga, ROI, dll        │              │   │
│  │  │    • outlet_packages → Paket & benefit      │              │   │
│  │  │    • ebooks → Judul, harga, deskripsi       │              │   │
│  │  │    • outlet_categories → Kategori            │              │   │
│  │  └────────────────────────────────────────────┘              │   │
│  │  ┌────────────────────────────────────────────┐              │   │
│  │  │ 4. CONVERSATION HISTORY (10 pesan terakhir) │              │   │
│  │  │    • Konteks percakapan sebelumnya          │              │   │
│  │  └────────────────────────────────────────────┘              │   │
│  │  ┌────────────────────────────────────────────┐              │   │
│  │  │ 5. USER MESSAGE (pesan terbaru)             │              │   │
│  │  │    • Pertanyaan/pesan dari user             │              │   │
│  │  └────────────────────────────────────────────┘              │   │
│  └──────────────────┬───────────────────────────────────────────┘   │
│                     │                                               │
└─────────────────────┼───────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     OPENAI API                                      │
│                                                                     │
│  Model: GPT-4o / GPT-5 (configurable via ai_config)                │
│  Temperature: 0.7 (configurable)                                    │
│  Max Tokens: 2048 (configurable)                                    │
│                                                                     │
│  Request → Proses AI → Response                                     │
└──────────────────────┬──────────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    DATABASE (PostgreSQL)                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─── Chat ───────────┐  ┌─── AI Knowledge ──────────────────┐     │
│  │ chat_conversations  │  │ ai_kb_categories                  │     │
│  │ • id                │  │ • id, name, slug                  │     │
│  │ • user_id           │  │ • description, sort_order         │     │
│  │ • title             │  │ • is_active                       │     │
│  │ • created_at        │  └───────────────────────────────────┘     │
│  │ • updated_at        │  ┌───────────────────────────────────┐     │
│  └─────────────────────┘  │ ai_knowledge_base                 │     │
│  ┌─────────────────────┐  │ • id, category_id                 │     │
│  │ chat_messages        │  │ • title, slug, content            │     │
│  │ • id                 │  │ • keywords[] (array)              │     │
│  │ • conversation_id    │  │ • priority, is_active             │     │
│  │ • user_id            │  └───────────────────────────────────┘     │
│  │ • role               │  ┌───────────────────────────────────┐     │
│  │ • content            │  │ ai_system_prompts                 │     │
│  │ • tokens_used        │  │ • id, name, prompt                │     │
│  │ • created_at         │  │ • is_active                       │     │
│  └─────────────────────┘  └───────────────────────────────────┘     │
│                            ┌───────────────────────────────────┐     │
│  ┌─── Business Data ──┐   │ ai_config                         │     │
│  │ outlets             │   │ • key: openai_model               │     │
│  │ outlet_packages     │   │ • key: openai_temperature         │     │
│  │ ebooks              │   │ • key: openai_max_tokens          │     │
│  │ outlet_categories   │   └───────────────────────────────────┘     │
│  └─────────────────────┘                                            │
└─────────────────────────────────────────────────────────────────────┘
```

## Alur Request Chat (Step by Step)

```
User Ketik Pesan
       │
       ▼
[ChatWidget.vue] ──POST /api/v1/mitra/chat──▶ [ChatHandler]
       │                                            │
       │                                     Validate JWT + Input
       │                                            │
       │                                            ▼
       │                                    [ChatService.Chat()]
       │                                            │
       │                               ┌────────────┼────────────┐
       │                               ▼            ▼            ▼
       │                        Get/Create    Load System    Load Knowledge
       │                        Conversation    Prompt         Base
       │                               │            │            │
       │                               ▼            ▼            ▼
       │                        Save User      Load Live     Build Chat
       │                        Message        Biz Data       History
       │                               │            │            │
       │                               └────────────┼────────────┘
       │                                            │
       │                                            ▼
       │                                    [buildGPTMessages()]
       │                                   System + KB + Data + History
       │                                            │
       │                                            ▼
       │                                    [OpenAI API Call]
       │                                     GPT-4o / GPT-5
       │                                            │
       │                                            ▼
       │                                    Save Assistant Message
       │                                    Generate Quick Actions
       │                                            │
       ◀────────── JSON Response ───────────────────┘
       │
       ▼
  Render Markdown
  Show Quick Actions
  Scroll to Bottom
```

## Alur Admin Mengelola AI

```
Admin Panel
       │
       ▼
[AISettings.vue] ──────▶ [AI Admin Handler]
       │                        │
       │              ┌────────┼────────┐───────────┐
       │              ▼        ▼        ▼           ▼
       │          Knowledge  System   AI Config   Cache
       │          Base CRUD  Prompts  (model,     Invalidate
       │                     CRUD     temp, etc)
       │              │        │        │           │
       │              ▼        ▼        ▼           ▼
       │          ai_knowledge ai_system ai_config  Service
       │          _base        _prompts             .InvalidateCache()
       │              │        │        │
       │              └────────┼────────┘
       │                       │
       │                       ▼
       │              Perubahan langsung efektif
       │              setelah cache expire (5 min)
       │              atau setelah invalidate
       ◀──────────────────────┘
```

## File Structure

```
backend/
├── cmd/api/main.go                              ← Entry point, wiring
├── config/config.go                             ← OpenAI config
├── .env                                         ← OPENAI_API_KEY
├── migrations/
│   ├── 034_create_chat_tables.up.sql            ← Chat tables
│   └── 035_create_ai_knowledge_base.up.sql      ← AI KB tables
├── internal/
│   ├── entity/chatbot.go                        ← All AI entities
│   ├── handler/
│   │   ├── chat_handler.go                      ← Mitra chat API
│   │   └── ai_admin_handler.go                  ← Admin AI settings API
│   ├── repository/postgres/
│   │   ├── chat_repo.go                         ← Chat CRUD
│   │   └── ai_knowledge_repo.go                 ← KB + Prompt + Config CRUD
│   ├── router/
│   │   ├── panel_routes.go                      ← Admin AI routes
│   │   └── mitra_routes.go                      ← Mitra chat routes
│   ├── seeder/ai_seeder.go                      ← KB seeder (25+ entries)
│   └── service/chatbot/
│       ├── service.go                           ← Core AI service
│       └── openai_client.go                     ← OpenAI API client

frontend/
├── mitra/src/
│   ├── components/ChatWidget.vue                ← Floating chat widget
│   ├── layouts/MitraLayout.vue                  ← Widget integrated
│   └── services/api.js                          ← Chat API functions
└── panel/src/
    ├── pages/ai/AISettings.vue                  ← Admin AI settings page
    ├── router/index.js                          ← AI route added
    ├── services/api.js                          ← AI admin API functions
    └── components/AppSidebar.vue                ← AI menu link
```

## Konfigurasi yang Bisa Diatur dari Admin Panel

| Konfigurasi    | Tabel               | Keterangan                 |
| -------------- | ------------------- | -------------------------- |
| System Prompt  | `ai_system_prompts` | Aturan & karakter AI       |
| Knowledge Base | `ai_knowledge_base` | Pengetahuan bisnis AI      |
| KB Categories  | `ai_kb_categories`  | Kategori knowledge         |
| Model AI       | `ai_config`         | gpt-4o, gpt-4o-mini, gpt-5 |
| Temperature    | `ai_config`         | 0.0 - 1.0 (kreativitas)    |
| Max Tokens     | `ai_config`         | Panjang respons maksimal   |

## Keamanan & Guardrails

1. **System Prompt yang ketat** — AI hanya boleh menjawab topik bisnis outlet
2. **JWT Authentication** — Setiap request chat harus login
3. **Rate limiting** — Mencegah abuse
4. **Token tracking** — Monitor penggunaan token per pesan
5. **Fallback response** — Jika OpenAI down, tampilkan pesan fallback
6. **Cache layer** — Mengurangi query DB berulang
7. **Input validation** — Maksimal 1000 karakter per pesan
