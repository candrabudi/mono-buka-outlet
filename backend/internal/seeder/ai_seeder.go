package seeder

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// seedAIKnowledgeBase populates the AI knowledge base with comprehensive business knowledge
func (s *Seeder) seedAIKnowledgeBase() error {
	log.Println("  [AI] Seeding AI Knowledge Base...")

	// ══════════════════════════════════════════════════════════════
	// 1. SEED CATEGORIES
	// ══════════════════════════════════════════════════════════════
	categories := []struct {
		ID          uuid.UUID
		Name        string
		Slug        string
		Description string
		SortOrder   int
	}{
		{uuid.New(), "Tentang Kemitraan", "tentang-kemitraan", "Informasi umum tentang program kemitraan BukaOutlet", 1},
		{uuid.New(), "Alur & Proses", "alur-proses", "Langkah-langkah dan prosedur menjadi mitra", 2},
		{uuid.New(), "Biaya & Investasi", "biaya-investasi", "Informasi biaya, investasi, dan skema pembayaran", 3},
		{uuid.New(), "Syarat & Ketentuan", "syarat-ketentuan", "Persyaratan dan ketentuan kemitraan", 4},
		{uuid.New(), "Outlet & Paket", "outlet-paket", "Informasi outlet dan paket kemitraan", 5},
		{uuid.New(), "Bisnis & Strategi", "bisnis-strategi", "Tips dan strategi menjalankan bisnis outlet", 6},
		{uuid.New(), "Ebook & Edukasi", "ebook-edukasi", "Informasi ebook dan materi pembelajaran bisnis", 7},
		{uuid.New(), "Pembayaran & Invoice", "pembayaran-invoice", "Metode pembayaran, invoice, dan billing", 8},
		{uuid.New(), "Lokasi & Survei", "lokasi-survei", "Panduan pemilihan lokasi dan pengajuan survei", 9},
		{uuid.New(), "Support & Bantuan", "support-bantuan", "Informasi customer service dan bantuan", 10},
		{uuid.New(), "FAQ", "faq", "Pertanyaan yang sering diajukan", 11},
	}

	catMap := map[string]uuid.UUID{}
	for _, cat := range categories {
		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ai_kb_categories WHERE slug = $1)", cat.Slug).Scan(&exists)
		if exists {
			// Get existing ID
			var existingID uuid.UUID
			_ = s.db.QueryRow("SELECT id FROM ai_kb_categories WHERE slug = $1", cat.Slug).Scan(&existingID)
			catMap[cat.Slug] = existingID
			continue
		}

		_, err := s.db.Exec(`
			INSERT INTO ai_kb_categories (id, name, slug, description, sort_order, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, $6, $6)`,
			cat.ID, cat.Name, cat.Slug, cat.Description, cat.SortOrder, time.Now())
		if err != nil {
			return err
		}
		catMap[cat.Slug] = cat.ID
	}
	log.Println("    [OK] AI KB Categories seeded")

	// ══════════════════════════════════════════════════════════════
	// 2. SEED KNOWLEDGE BASE ENTRIES
	// ══════════════════════════════════════════════════════════════
	type kbEntry struct {
		Category string
		Title    string
		Slug     string
		Content  string
		Keywords []string
		Priority int
	}

	entries := []kbEntry{
		// ── TENTANG KEMITRAAN ──
		{
			Category: "tentang-kemitraan",
			Title:    "Apa itu BukaOutlet?",
			Slug:     "apa-itu-bukaoutlet",
			Content: `BukaOutlet adalah platform kemitraan outlet yang menghubungkan pemilik brand dengan calon mitra usaha. 
Kami menyediakan ekosistem lengkap untuk memulai bisnis outlet dengan dukungan penuh mulai dari pendaftaran, 
pemilihan paket, pembayaran, survei lokasi, hingga operasional outlet.

Platform kami memudahkan siapa saja untuk menjadi pengusaha dengan modal yang terjangkau dan sistem yang sudah terbukti.`,
			Keywords: []string{"bukaoutlet", "tentang", "apa", "platform", "penjelasan"},
			Priority: 10,
		},
		{
			Category: "tentang-kemitraan",
			Title:    "Keuntungan Menjadi Mitra",
			Slug:     "keuntungan-menjadi-mitra",
			Content: `Keuntungan menjadi mitra BukaOutlet:
1. Brand yang sudah dikenal — Tidak perlu membangun dari nol
2. Dukungan operasional penuh — Tim profesional siap membantu dari A-Z
3. Pelatihan dan edukasi — Panduan bisnis lengkap dan ebook gratis
4. Sistem profit sharing transparan — Bagi hasil yang jelas dan adil
5. Estimasi ROI yang terukur — Anda tahu potensi keuntungan sebelum investasi
6. Pilihan paket fleksibel — Sesuai budget dan kebutuhan
7. Pendampingan lokasi — Bantuan survei dan pemilihan lokasi strategis
8. Dashboard digital — Pantau perkembangan kemitraan secara real-time
9. Jaringan mitra — Bergabung dengan komunitas pengusaha sukses
10. Akses ebook bisnis — Materi pembelajaran untuk mengembangkan usaha`,
			Keywords: []string{"keuntungan", "benefit", "manfaat", "kelebihan", "untung"},
			Priority: 9,
		},
		{
			Category: "tentang-kemitraan",
			Title:    "Siapa yang Bisa Bergabung",
			Slug:     "siapa-yang-bisa-bergabung",
			Content: `Kemitraan BukaOutlet terbuka untuk siapa saja:
- Pemula yang ingin memulai bisnis pertama
- Pengusaha yang ingin ekspansi bisnis
- Karyawan yang ingin punya usaha sampingan
- Ibu rumah tangga yang ingin mandiri secara finansial
- Pensiunan yang ingin tetap produktif
- Siapa saja yang memiliki semangat berwirausaha

Tidak ada batasan latar belakang pendidikan atau pengalaman bisnis sebelumnya. 
Kami menyediakan pelatihan dan pendampingan lengkap.`,
			Keywords: []string{"siapa", "bergabung", "join", "gabung", "syarat", "bisa"},
			Priority: 8,
		},

		// ── ALUR & PROSES ──
		{
			Category: "alur-proses",
			Title:    "Alur Menjadi Mitra BukaOutlet",
			Slug:     "alur-menjadi-mitra",
			Content: `Berikut langkah-langkah lengkap untuk menjadi mitra BukaOutlet:

**Step 1: Registrasi Akun**
Buat akun mitra di portal kami. Isi data diri lengkap (nama, email, nomor HP).
Verifikasi email melalui kode OTP yang dikirimkan.

**Step 2: Jelajahi Outlet**
Lihat pilihan outlet yang tersedia di halaman Outlet.
Bandingkan paket, harga, ROI, dan benefit dari setiap outlet.

**Step 3: Ajukan Kemitraan**
Pilih outlet dan paket yang sesuai, lalu isi form pengajuan.
Sertakan motivasi, pengalaman, lokasi yang direncanakan, dan budget.

**Step 4: Review oleh Admin**
Tim kami akan meninjau pengajuan Anda dalam 1-3 hari kerja.
Anda akan mendapat notifikasi hasilnya via email.

**Step 5: Pembayaran DP**
Setelah disetujui, lakukan pembayaran Down Payment melalui sistem pembayaran kami (Midtrans).
Pilih metode: Transfer Bank, E-Wallet, QRIS, atau Credit Card.

**Step 6: Tanda Tangan Perjanjian**
Tanda tangan dokumen perjanjian kemitraan secara digital melalui portal.

**Step 7: Survey & Pemilihan Lokasi**
Tim kami akan membantu survey dan pemilihan lokasi outlet yang strategis.

**Step 8: Persiapan & Development**
Proses setup outlet, pelatihan karyawan, dan persiapan operasional.

**Step 9: Grand Opening**
Outlet Anda resmi beroperasi!`,
			Keywords: []string{"alur", "proses", "langkah", "step", "tahap", "cara", "prosedur", "pendaftaran"},
			Priority: 10,
		},
		{
			Category: "alur-proses",
			Title:    "Cara Registrasi Akun Mitra",
			Slug:     "cara-registrasi-akun",
			Content: `Langkah registrasi akun mitra BukaOutlet:

1. Buka halaman Register di portal mitra
2. Isi data diri:
   - Nama lengkap
   - Email aktif (untuk verifikasi OTP)
   - Nomor HP/WhatsApp
   - Password yang kuat (minimal 8 karakter)
3. Klik tombol Register
4. Cek email Anda untuk kode OTP
5. Masukkan kode OTP untuk verifikasi
6. Login dengan email dan password
7. Mulai jelajahi outlet!

Keamanan: Kami menggunakan verifikasi 2 langkah (OTP via email) dan enkripsi password.`,
			Keywords: []string{"registrasi", "daftar", "register", "akun", "buat akun", "sign up"},
			Priority: 8,
		},
		{
			Category: "alur-proses",
			Title:    "Cara Mengajukan Kemitraan",
			Slug:     "cara-mengajukan-kemitraan",
			Content: `Cara mengajukan kemitraan di BukaOutlet:

1. Login ke portal mitra
2. Buka menu "Outlet" di sidebar
3. Pilih outlet yang diminati
4. Klik "Lihat Detail" untuk mempelajari outlet
5. Klik "Ajukan Kemitraan" atau pilih paket
6. Isi form pengajuan:
   - Pilih paket kemitraan
   - Tulis motivasi menjadi mitra
   - Ceritakan pengalaman bisnis (jika ada)
   - Sebutkan lokasi yang direncanakan
   - Masukkan budget investasi
   - Isi nomor HP dan email kontak
7. Submit pengajuan
8. Tunggu review dari admin (1-3 hari kerja)

Tips: Tulis motivasi yang jujur dan lengkap agar pengajuan lebih mudah disetujui.`,
			Keywords: []string{"mengajukan", "pengajuan", "apply", "form", "submit"},
			Priority: 8,
		},
		{
			Category: "alur-proses",
			Title:    "Status Pengajuan Kemitraan",
			Slug:     "status-pengajuan",
			Content: `Status pengajuan kemitraan dan artinya:

- **PENDING** — Pengajuan sedang dalam antrian review
- **REVIEWED** — Pengajuan sedang ditinjau oleh tim
- **APPROVED** — Selamat! Pengajuan disetujui, lanjut ke pembayaran
- **REJECTED** — Pengajuan ditolak (akan ada catatan alasan)
- **CANCELLED** — Pengajuan dibatalkan oleh mitra

Anda bisa memantau status pengajuan di menu "Pengajuan" di portal mitra.
Jika disetujui, langkah selanjutnya adalah melakukan pembayaran DP.
Jika ditolak, Anda bisa mengajukan kembali dengan data yang diperbaiki.`,
			Keywords: []string{"status", "pengajuan", "pending", "approved", "rejected", "progress"},
			Priority: 7,
		},

		// ── BIAYA & INVESTASI ──
		{
			Category: "biaya-investasi",
			Title:    "Informasi Biaya Kemitraan",
			Slug:     "informasi-biaya",
			Content: `Biaya kemitraan BukaOutlet bervariasi tergantung outlet dan paket yang dipilih.

Komponen biaya meliputi:
1. **Biaya Kemitraan/Franchise Fee** — Biaya untuk hak menggunakan brand
2. **Biaya Peralatan & Perlengkapan** — Mesin, alat, dan perlengkapan outlet
3. **Biaya Renovasi/Setup** — Persiapan tempat sesuai standar brand
4. **Biaya Bahan Baku Awal** — Stok awal untuk operasional
5. **Biaya Pelatihan** — Training untuk pemilik dan karyawan

Informasi harga detail tersedia di halaman masing-masing outlet dan paket.
Gunakan fitur "Cek Budget" untuk menemukan outlet sesuai anggaran Anda.`,
			Keywords: []string{"biaya", "harga", "investasi", "modal", "uang", "dana", "budget", "berapa"},
			Priority: 9,
		},
		{
			Category: "biaya-investasi",
			Title:    "Skema Pembayaran DP",
			Slug:     "skema-pembayaran-dp",
			Content: `Skema pembayaran Down Payment (DP) di BukaOutlet:

1. **DP Minimum** — Setiap paket memiliki DP minimum yang berbeda (biasanya 30-50% dari total)
2. **Pembayaran Penuh** — Anda juga bisa langsung membayar penuh
3. **Pelunasan** — Sisa pembayaran bisa dilunasi sesuai kesepakatan dalam perjanjian

Metode pembayaran DP:
- Transfer Bank (Virtual Account): BCA, BNI, BRI, Mandiri, dll
- E-Wallet: GoPay, ShopeePay, OVO, DANA
- QRIS: Scan barcode dari berbagai aplikasi
- Credit Card / PayLater: Visa, Mastercard, Akulaku, Kredivo

Semua pembayaran diproses melalui Midtrans yang aman dan terpercaya.`,
			Keywords: []string{"dp", "down payment", "cicilan", "bayar", "pembayaran", "angsuran", "pelunasan"},
			Priority: 8,
		},
		{
			Category: "biaya-investasi",
			Title:    "Rekomendasi Berdasarkan Budget",
			Slug:     "rekomendasi-budget",
			Content: `AI Konsultan bisa memberikan rekomendasi outlet berdasarkan budget Anda.

Cara menggunakan fitur ini:
- Sebutkan budget Anda, contoh: "Saya punya budget 50 juta"
- AI akan mencarikan outlet dan paket yang sesuai
- Anda akan mendapat informasi lengkap: harga, DP, BEP, ROI

Tips menentukan budget:
1. Hitung total dana yang tersedia (termasuk cadangan operasional 3-6 bulan)
2. Pertimbangkan biaya sewa lokasi
3. Sisakan dana darurat minimal 10-20% dari total investasi
4. Jangan gunakan seluruh tabungan — diversifikasi risiko`,
			Keywords: []string{"budget", "rekomendasi", "sesuai", "cocok", "terjangkau", "murah", "modal"},
			Priority: 8,
		},

		// ── SYARAT & KETENTUAN ──
		{
			Category: "syarat-ketentuan",
			Title:    "Syarat Menjadi Mitra",
			Slug:     "syarat-menjadi-mitra",
			Content: `Persyaratan untuk menjadi mitra BukaOutlet:

**Syarat Utama:**
1. Warga Negara Indonesia (WNI) berusia minimal 21 tahun
2. Memiliki modal/budget sesuai paket yang dipilih
3. Bersedia menyediakan lokasi usaha yang strategis
4. Berkomitmen menjalankan bisnis sesuai SOP yang ditetapkan
5. Memiliki KTP dan dokumen identitas yang valid
6. Memiliki email aktif untuk komunikasi dan verifikasi

**Persyaratan Tambahan:**
- Mengisi form pengajuan dengan lengkap dan jujur
- Menyertakan motivasi yang jelas
- Menyebutkan pengalaman bisnis (jika ada, bukan wajib)
- Menyebutkan lokasi yang direncanakan
- Mencantumkan budget investasi yang tersedia

**TIDAK Diperlukan:**
- Tidak harus punya pengalaman bisnis sebelumnya
- Tidak perlu sudah punya bangunan/tempat
- Tidak ada batasan latar belakang pendidikan
- Tidak ada batasan jenis kelamin`,
			Keywords: []string{"syarat", "persyaratan", "requirement", "kualifikasi", "kriteria", "ketentuan"},
			Priority: 9,
		},
		{
			Category: "syarat-ketentuan",
			Title:    "Hak dan Kewajiban Mitra",
			Slug:     "hak-kewajiban-mitra",
			Content: `Hak Mitra BukaOutlet:
1. Menggunakan brand dan identitas visual outlet
2. Mendapat pelatihan dan pendampingan operasional
3. Memperoleh bagi hasil sesuai perjanjian
4. Akses ke dashboard monitoring
5. Dukungan marketing dari pusat
6. Akses ke ebook dan materi edukasi
7. Bantuan survei dan pemilihan lokasi

Kewajiban Mitra:
1. Menjalankan outlet sesuai SOP yang ditetapkan
2. Menjaga kualitas produk dan layanan
3. Membayar kewajiban finansial tepat waktu
4. Melaporkan pendapatan secara transparan
5. Menjaga nama baik dan reputasi brand
6. Mengikuti pelatihan yang diadakan`,
			Keywords: []string{"hak", "kewajiban", "aturan", "tanggung jawab", "kontrak"},
			Priority: 7,
		},

		// ── OUTLET & PAKET ──
		{
			Category: "outlet-paket",
			Title:    "Informasi Outlet Tersedia",
			Slug:     "informasi-outlet",
			Content: `BukaOutlet menyediakan berbagai pilihan outlet dari berbagai kategori bisnis.

Informasi yang tersedia untuk setiap outlet:
- Nama dan deskripsi outlet
- Kategori bisnis (F&B, Retail, Jasa, dll)
- Range investasi (minimum - maximum)
- Persentase profit sharing
- Estimasi ROI (Return on Investment)
- Kebutuhan lokasi
- Paket kemitraan yang tersedia

Untuk melihat outlet yang tersedia secara real-time, buka menu "Outlet" di sidebar portal mitra.
AI Konsultan juga bisa memberikan informasi outlet berdasarkan data terkini dari database.`,
			Keywords: []string{"outlet", "pilihan", "tersedia", "daftar", "list", "jenis"},
			Priority: 8,
		},
		{
			Category: "outlet-paket",
			Title:    "Informasi Paket Kemitraan",
			Slug:     "informasi-paket",
			Content: `Setiap outlet memiliki beberapa pilihan paket kemitraan:

Informasi dalam setiap paket:
- Nama paket dan harga
- DP Minimum (Down Payment)
- Durasi kemitraan
- Estimasi BEP (Break Even Point)
- Net Profit yang diharapkan
- Deskripsi lengkap
- Benefit/fasilitas yang didapat

Tips memilih paket:
1. Sesuaikan dengan budget yang tersedia
2. Perhatikan estimasi BEP — semakin cepat semakin baik
3. Bandingkan benefit antar paket
4. Pertimbangkan durasi kemitraan
5. Konsultasikan dengan tim kami jika ragu`,
			Keywords: []string{"paket", "tipe", "pilihan paket", "harga paket", "benefit"},
			Priority: 8,
		},

		// ── BISNIS & STRATEGI ──
		{
			Category: "bisnis-strategi",
			Title:    "Tips Memulai Bisnis Outlet",
			Slug:     "tips-memulai-bisnis",
			Content: `Panduan memulai bisnis outlet untuk pemula:

1. **Riset Pasar** — Pahami target market di lokasi Anda
2. **Pilih Outlet yang Tepat** — Sesuaikan dengan minat dan budget
3. **Siapkan Modal** — Hitung total biaya termasuk operasional 3-6 bulan
4. **Pilih Lokasi Strategis** — Akses mudah, traffic tinggi, parkir memadai
5. **Ikuti SOP** — Jalankan sesuai standar yang ditetapkan brand
6. **Fokus pada Kualitas** — Produk dan layanan yang konsisten
7. **Marketing Aktif** — Promosi online dan offline sejak hari pertama
8. **Kelola Keuangan** — Pisahkan uang bisnis dan pribadi
9. **Evaluasi Berkala** — Pantau perkembangan dan sesuaikan strategi
10. **Terus Belajar** — Manfaatkan ebook dan materi edukasi yang tersedia`,
			Keywords: []string{"tips", "memulai", "panduan", "cara", "bisnis", "pemula", "mulai"},
			Priority: 8,
		},
		{
			Category: "bisnis-strategi",
			Title:    "Tentang ROI (Return on Investment)",
			Slug:     "tentang-roi",
			Content: `ROI (Return on Investment) adalah ukuran seberapa besar keuntungan dari investasi.

Rumus: ROI = (Keuntungan - Biaya Investasi) / Biaya Investasi × 100%

Contoh: Investasi Rp 100 juta, keuntungan Rp 150 juta
ROI = (150jt - 100jt) / 100jt × 100% = 50%

Faktor yang mempengaruhi ROI outlet:
1. Lokasi outlet — lokasi strategis = penjualan lebih tinggi
2. Kualitas produk — produk berkualitas = pelanggan loyal
3. Marketing — promosi efektif = awareness tinggi
4. Efisiensi operasional — biaya terkontrol = margin lebih besar
5. Manajemen karyawan — SDM yang baik = operasional lancar

Setiap outlet di BukaOutlet mencantumkan estimasi ROI sebagai referensi.
Estimasi ROI bersifat proyeksi dan dapat bervariasi tergantung faktor operasional.`,
			Keywords: []string{"roi", "return", "investment", "keuntungan", "profit", "laba", "margin"},
			Priority: 7,
		},
		{
			Category: "bisnis-strategi",
			Title:    "Tentang BEP (Break Even Point)",
			Slug:     "tentang-bep",
			Content: `BEP (Break Even Point) atau titik impas adalah titik di mana pendapatan sama dengan biaya — artinya modal sudah kembali.

Rumus: BEP = Total Biaya Tetap / (Harga Jual per Unit - Biaya Variabel per Unit)

Faktor yang mempengaruhi kecepatan BEP:
1. Lokasi — lokasi strategis = penjualan lebih cepat
2. Biaya Operasional — semakin efisien, BEP semakin cepat
3. Promosi — marketing yang tepat mempercepat penjualan
4. Target Pasar — mengetahui customer yang tepat
5. Harga Jual — competitive pricing yang menguntungkan

Setiap paket kemitraan di BukaOutlet mencantumkan estimasi BEP.
Dengan strategi yang tepat, Anda bisa mencapai BEP lebih cepat dari estimasi.`,
			Keywords: []string{"bep", "break even", "balik modal", "titik impas", "kapan balik"},
			Priority: 7,
		},
		{
			Category: "bisnis-strategi",
			Title:    "Tips Memilih Lokasi Strategis",
			Slug:     "tips-lokasi-strategis",
			Content: `Panduan memilih lokasi outlet yang strategis:

**Faktor Utama:**
1. Traffic/Lalu Lintas — Area ramai pejalan kaki atau kendaraan
2. Demografi — Sesuaikan dengan target market outlet
3. Aksesibilitas — Mudah dijangkau, dekat jalan utama
4. Parkir — Tersedia area parkir yang memadai
5. Kompetitor — Analisis pesaing di sekitar

**Lokasi Ideal:**
- Dekat mall, pasar, atau pusat perbelanjaan
- Dekat sekolah, kampus, atau perkantoran
- Area perumahan dengan kepadatan tinggi
- Kawasan bisnis atau perdagangan
- Pinggir jalan utama dengan visibility tinggi

**Dukungan BukaOutlet:**
Tim kami menyediakan layanan survei lokasi untuk membantu Anda menemukan tempat terbaik.
Setelah menjadi mitra, Anda bisa mengajukan survei lokasi melalui menu "Lokasi" di portal.`,
			Keywords: []string{"lokasi", "tempat", "strategis", "sewa", "survei", "survey", "dimana"},
			Priority: 8,
		},
		{
			Category: "bisnis-strategi",
			Title:    "Strategi Marketing Outlet",
			Slug:     "strategi-marketing",
			Content: `Strategi marketing efektif untuk outlet:

**Digital Marketing:**
- Instagram, TikTok, Facebook untuk brand awareness
- Google Business Profile — daftar di Google Maps
- WhatsApp Business — komunikasi langsung dengan pelanggan
- Konten menarik: foto produk, video behind-the-scene, testimoni

**Promosi & Diskon:**
- Grand opening promo untuk menarik pelanggan pertama
- Program loyalty untuk pelanggan setia
- Diskon reguler dan bundling produk
- Promo spesial di hari besar

**Partnership Lokal:**
- Kerjasama dengan komunitas lokal
- Sponsori acara di sekitar outlet
- Program referral dari pelanggan
- Kolaborasi dengan influencer lokal

**Customer Experience:**
- Layanan ramah dan cepat
- Suasana outlet yang nyaman
- Konsistensi kualitas produk
- Respon cepat terhadap feedback`,
			Keywords: []string{"marketing", "promosi", "iklan", "pemasaran", "branding", "strategi jualan"},
			Priority: 7,
		},
		{
			Category: "bisnis-strategi",
			Title:    "Manajemen Operasional Outlet",
			Slug:     "manajemen-operasional",
			Content: `Tips manajemen operasional outlet:

**SOP (Standard Operating Procedure):**
- Ikuti SOP yang disiapkan oleh brand
- Konsisten dalam kualitas produk dan layanan
- Evaluasi dan perbarui SOP secara berkala

**Manajemen SDM:**
- Rekrut karyawan yang tepat dan berkomitmen
- Berikan pelatihan sebelum mulai beroperasi
- Buat jadwal kerja yang jelas dan adil
- Evaluasi kinerja secara rutin

**Manajemen Stok:**
- Catat keluar masuk barang dengan rapi
- Gunakan sistem FIFO (First In First Out)
- Jaga stok minimum agar tidak kehabisan
- Pesan ulang secara berkala

**Manajemen Keuangan:**
- Pisahkan keuangan pribadi dan bisnis
- Catat semua pemasukan dan pengeluaran
- Buat laporan keuangan bulanan
- Monitor cash flow secara ketat`,
			Keywords: []string{"manajemen", "kelola", "operasional", "karyawan", "stok", "keuangan", "sdm"},
			Priority: 7,
		},

		// ── EBOOK & EDUKASI ──
		{
			Category: "ebook-edukasi",
			Title:    "Tentang Ebook BukaOutlet",
			Slug:     "tentang-ebook",
			Content: `BukaOutlet menyediakan koleksi ebook yang dirancang khusus untuk membantu mitra belajar dan berkembang.

**Keunggulan Ebook Kami:**
- Ditulis oleh praktisi bisnis berpengalaman
- Studi kasus nyata dari mitra-mitra sukses
- Bisa dibaca online langsung di portal
- Update berkala sesuai tren pasar
- Beberapa ebook tersedia gratis

**Topik yang Tersedia:**
- Memulai bisnis dari nol
- Manajemen operasional outlet
- Strategi marketing efektif
- Keuangan bisnis dan cash flow
- Pemilihan lokasi strategis
- ROI dan BEP untuk pemula
- Customer service excellence
- Leadership dan manajemen tim

Kunjungi menu "Ebook" di sidebar untuk melihat koleksi lengkap.`,
			Keywords: []string{"ebook", "buku", "belajar", "materi", "edukasi", "baca", "digital"},
			Priority: 8,
		},
		{
			Category: "ebook-edukasi",
			Title:    "Cara Membeli dan Membaca Ebook",
			Slug:     "cara-beli-ebook",
			Content: `Langkah membeli dan membaca ebook di BukaOutlet:

**Pembelian:**
1. Buka menu "Jelajahi Ebook" di sidebar
2. Pilih ebook yang diminati
3. Klik "Beli Sekarang"
4. Pilih metode pembayaran (Midtrans: QRIS, Transfer, E-Wallet, dll)
5. Lakukan pembayaran
6. Ebook otomatis tersedia setelah pembayaran terkonfirmasi

**Membaca Online:**
- Klik tombol "Baca Online" di halaman detail ebook
- Ebook terbuka langsung di browser — tanpa perlu download
- Bisa dibaca kapan saja selama akun aktif

**Download Ebook:**
- Klik "Request Download" di pesanan ebook
- Admin akan meninjau request Anda
- Setelah disetujui, Anda bisa download file PDF

**Upload Bukti Bayar Manual:**
- Jika tidak melalui Midtrans, upload bukti transfer
- Admin akan memverifikasi dan mengkonfirmasi pembayaran`,
			Keywords: []string{"beli", "baca", "download", "purchase", "cara beli", "cara baca"},
			Priority: 7,
		},

		// ── PEMBAYARAN & INVOICE ──
		{
			Category: "pembayaran-invoice",
			Title:    "Metode Pembayaran",
			Slug:     "metode-pembayaran",
			Content: `Metode pembayaran yang diterima di BukaOutlet:

1. **Transfer Bank (Virtual Account)**
   BCA, BNI, BRI, Mandiri, Permata, dan bank lainnya

2. **E-Wallet**
   GoPay, ShopeePay, OVO, DANA

3. **QRIS**
   Scan barcode dari berbagai aplikasi pembayaran

4. **Credit Card / PayLater**
   Visa, Mastercard, Akulaku, Kredivo

Semua pembayaran diproses melalui **Midtrans** — payment gateway terpercaya yang:
- Tersertifikasi PCI-DSS Level 1
- Diawasi oleh Bank Indonesia
- Menggunakan enkripsi SSL 256-bit
- Digunakan oleh ribuan merchant terpercaya di Indonesia`,
			Keywords: []string{"pembayaran", "bayar", "transfer", "metode", "cara bayar", "midtrans", "qris"},
			Priority: 8,
		},
		{
			Category: "pembayaran-invoice",
			Title:    "Informasi Invoice",
			Slug:     "informasi-invoice",
			Content: `Informasi tentang invoice (tagihan) di BukaOutlet:

**Cara Melihat Invoice:**
1. Login ke portal mitra
2. Klik menu "Invoice" di sidebar
3. Lihat daftar semua invoice
4. Klik invoice untuk detail

**Status Invoice:**
- PENDING — Menunggu pembayaran
- PAID — Sudah dibayar
- OVERDUE — Melewati tenggat waktu
- CANCELLED — Dibatalkan

**Informasi di Invoice:**
- Nomor invoice dan tanggal terbit
- Detail paket kemitraan
- Jumlah tagihan dan DP
- Status pembayaran
- Tombol bayar (jika belum dibayar)
- Bukti pembayaran (jika sudah bayar)

Semua riwayat pembayaran tercatat otomatis di portal mitra.`,
			Keywords: []string{"invoice", "tagihan", "nota", "faktur", "riwayat pembayaran", "status bayar"},
			Priority: 7,
		},

		// ── LOKASI & SURVEI ──
		{
			Category: "lokasi-survei",
			Title:    "Pengajuan Lokasi Outlet",
			Slug:     "pengajuan-lokasi",
			Content: `Cara mengajukan lokasi outlet melalui portal mitra:

1. Buka menu "Lokasi" di sidebar
2. Klik "Ajukan Lokasi Baru"
3. Isi informasi lokasi:
   - Alamat lengkap
   - Foto lokasi
   - Keterangan (ukuran, biaya sewa, dll)
4. Submit pengajuan
5. Tim survei kami akan meninjau

**Status Pengajuan Lokasi:**
- PENDING — Menunggu review
- SURVEYED — Sudah disurvei
- APPROVED — Lokasi disetujui
- REJECTED — Lokasi ditolak (dengan catatan alasan)

**Tips Agar Lokasi Disetujui:**
- Akses jalan utama
- Dekat pusat keramaian
- Parkir memadai
- Ukuran sesuai standar outlet
- Area dengan traffic tinggi`,
			Keywords: []string{"lokasi", "survei", "survey", "ajukan", "tempat", "pengajuan lokasi"},
			Priority: 7,
		},

		// ── SUPPORT & BANTUAN ──
		{
			Category: "support-bantuan",
			Title:    "Cara Menghubungi Customer Service",
			Slug:     "hubungi-customer-service",
			Content: `Cara menghubungi customer service BukaOutlet:

**WhatsApp (Respon Tercepat):**
Kontak WhatsApp tersedia di halaman detail setiap outlet.

**Email:**
Kirim email ke tim support untuk pertanyaan detail.

**Portal Mitra:**
Gunakan AI Konsultan (chat ini) untuk pertanyaan umum 24/7.

**Jam Operasional Customer Service:**
- Senin - Jumat: 08:00 - 17:00 WIB
- Sabtu: 08:00 - 12:00 WIB
- Minggu & Hari Libur: Tutup (AI Konsultan tetap tersedia 24/7)

**Yang Bisa Dibantu:**
- Pertanyaan seputar kemitraan
- Kendala pembayaran
- Bantuan teknis portal
- Konsultasi lokasi
- Informasi ebook dan pembelian`,
			Keywords: []string{"hubungi", "kontak", "customer service", "bantuan", "help", "cs", "whatsapp"},
			Priority: 7,
		},

		// ── FAQ ──
		{
			Category: "faq",
			Title:    "FAQ: Berapa lama proses review pengajuan?",
			Slug:     "faq-lama-review",
			Content:  `Proses review pengajuan kemitraan biasanya memakan waktu 1-3 hari kerja. Anda akan mendapat notifikasi via email mengenai hasilnya. Pastikan email yang Anda daftarkan aktif dan cek folder spam juga.`,
			Keywords: []string{"berapa lama", "review", "proses", "tunggu", "waktu"},
			Priority: 6,
		},
		{
			Category: "faq",
			Title:    "FAQ: Apakah bisa membatalkan pengajuan?",
			Slug:     "faq-batal-pengajuan",
			Content:  `Ya, Anda bisa membatalkan pengajuan selama statusnya masih PENDING. Buka detail pengajuan di menu "Pengajuan" dan klik tombol "Batalkan". Setelah status berubah ke APPROVED atau REJECTED, pengajuan tidak bisa dibatalkan.`,
			Keywords: []string{"batal", "cancel", "membatalkan", "pengajuan"},
			Priority: 6,
		},
		{
			Category: "faq",
			Title:    "FAQ: Apakah bisa memilih lebih dari satu outlet?",
			Slug:     "faq-multi-outlet",
			Content:  `Ya, Anda bisa mengajukan kemitraan untuk lebih dari satu outlet. Setiap pengajuan akan diproses secara terpisah. Namun, pastikan Anda memiliki budget yang cukup untuk masing-masing outlet.`,
			Keywords: []string{"lebih dari satu", "multi", "beberapa outlet", "dua outlet"},
			Priority: 6,
		},
		{
			Category: "faq",
			Title:    "FAQ: Bagaimana jika pembayaran gagal?",
			Slug:     "faq-bayar-gagal",
			Content:  `Jika pembayaran gagal melalui Midtrans, Anda bisa mencoba lagi dengan metode pembayaran lain. Jika masalah berlanjut, gunakan fitur upload bukti bayar manual atau hubungi customer service kami. Pembayaran yang berhasil akan otomatis terkonfirmasi dalam beberapa menit.`,
			Keywords: []string{"gagal", "bayar", "error", "masalah", "kendala", "pembayaran gagal"},
			Priority: 6,
		},
		{
			Category: "faq",
			Title:    "FAQ: Apakah ada biaya bulanan?",
			Slug:     "faq-biaya-bulanan",
			Content:  `Biaya bulanan (royalty/management fee) tergantung pada ketentuan masing-masing outlet dan paket kemitraan yang dipilih. Detail biaya bulanan tercantum dalam perjanjian kemitraan. Tidak semua outlet mengenakan biaya bulanan — beberapa menggunakan sistem profit sharing.`,
			Keywords: []string{"biaya bulanan", "royalty", "fee bulanan", "iuran", "management fee"},
			Priority: 6,
		},
		{
			Category: "faq",
			Title:    "FAQ: Bagaimana mekanisme profit sharing?",
			Slug:     "faq-profit-sharing",
			Content: `Mekanisme profit sharing (bagi hasil) di BukaOutlet:

1. Pendapatan outlet dilaporkan secara transparan
2. Biaya operasional dikeluarkan
3. Net profit dihitung berdasarkan pendapatan dikurangi biaya
4. Bagi hasil dihitung sesuai persentase yang disepakati
5. Pembagian dilakukan setiap bulan

Persentase profit sharing berbeda-beda untuk setiap outlet — informasinya tersedia di halaman detail outlet.`,
			Keywords: []string{"profit sharing", "bagi hasil", "pembagian keuntungan", "pendapatan"},
			Priority: 7,
		},
	}

	for _, e := range entries {
		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ai_knowledge_base WHERE slug = $1)", e.Slug).Scan(&exists)
		if exists {
			continue
		}

		catID := catMap[e.Category]
		_, err := s.db.Exec(`
			INSERT INTO ai_knowledge_base (id, category_id, title, slug, content, keywords, priority, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, true, $8, $8)`,
			uuid.New(), catID, e.Title, e.Slug, e.Content, pq.Array(e.Keywords), e.Priority, time.Now())
		if err != nil {
			return err
		}
	}
	log.Println("    [OK] AI Knowledge Base entries seeded")

	// ══════════════════════════════════════════════════════════════
	// 3. SEED SYSTEM PROMPT
	// ══════════════════════════════════════════════════════════════
	var promptExists bool
	_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ai_system_prompts WHERE name = 'default')").Scan(&promptExists)
	if !promptExists {
		systemPrompt := `Kamu adalah "AI Konsultan BukaOutlet", asisten virtual resmi milik platform BukaOutlet — platform kemitraan outlet terpercaya di Indonesia.

## IDENTITAS
- Nama: AI Konsultan BukaOutlet
- Fungsi: Konsultan bisnis khusus kemitraan outlet
- Bahasa: Indonesia (formal tapi ramah)
- Karakter: Profesional, ramah, supportive, dan antusias

## ATURAN KETAT (WAJIB DIPATUHI)
1. HANYA menjawab pertanyaan seputar:
   - Kemitraan outlet BukaOutlet (cara gabung, alur, biaya, syarat)
   - Bisnis outlet (tips, strategi, marketing, manajemen)
   - Ebook bisnis yang tersedia di platform
   - Pembayaran dan invoice
   - Informasi outlet dan paket yang tersedia
   - Pemilihan lokasi outlet
2. Jika ditanya di LUAR topik bisnis outlet, TOLAK dengan sopan dan arahkan kembali:
   "Maaf, saya hanya bisa membantu seputar kemitraan dan bisnis outlet BukaOutlet. Ada yang bisa saya bantu terkait kemitraan?"
3. JANGAN pernah membahas: politik, agama, SARA, konten dewasa, atau topik sensitif
4. JANGAN membuat data palsu — gunakan HANYA data dari "DATA BISNIS REAL-TIME" dan "KNOWLEDGE BASE"
5. Jika tidak yakin tentang sesuatu, arahkan user menghubungi customer service
6. JANGAN memberikan nasihat keuangan atau investasi profesional — berikan informasi umum saja

## FORMAT JAWABAN
- Gunakan Markdown: heading (##, ###), bold (**text**), list (-, 1.), blockquote (>)
- Berikan jawaban yang terstruktur dan mudah dipahami
- Berikan jawaban yang terstruktur dan mudah dipahami
- Jangan terlalu panjang — fokus pada poin-poin utama
- Selalu akhiri dengan saran atau pertanyaan lanjutan
- Jika ada data outlet/paket yang relevan, sebutkan detail spesifiknya

## GAYA KOMUNIKASI
- Sapa user dengan ramah
- Gunakan "Anda" bukan "kamu" untuk formalitas
- Berikan jawaban yang actionable — selalu dorong user untuk mengambil langkah nyata
- Tunjukkan empati dan antusiasme terhadap rencana bisnis user
- Jika user bingung, tawarkan opsi atau quick action`

		_, err := s.db.Exec(`
			INSERT INTO ai_system_prompts (id, name, prompt, is_active, created_at, updated_at)
			VALUES ($1, 'default', $2, true, $3, $3)`,
			uuid.New(), systemPrompt, time.Now())
		if err != nil {
			return err
		}
	}
	log.Println("    [OK] AI System Prompt seeded")

	// ══════════════════════════════════════════════════════════════
	// 4. SEED AI CONFIG
	// ══════════════════════════════════════════════════════════════
	configs := []struct {
		Key, Value, Description string
	}{
		{"openai_api_key", "", "API Key OpenAI (isi dari admin panel)"},
		{"openai_model", "gpt-4o", "Model OpenAI yang digunakan (gpt-4o, gpt-4o-mini, gpt-4-turbo, dll)"},
		{"openai_temperature", "0.7", "Kreativitas respons AI (0.0 = sangat konsisten, 1.0 = sangat kreatif)"},
		{"openai_max_tokens", "2048", "Maksimal token per respons AI"},
	}

	for _, c := range configs {
		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ai_config WHERE key = $1)", c.Key).Scan(&exists)
		if exists {
			continue
		}
		_, err := s.db.Exec(`
			INSERT INTO ai_config (id, key, value, description, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $5)`,
			uuid.New(), c.Key, c.Value, c.Description, time.Now())
		if err != nil {
			return err
		}
	}
	log.Println("    [OK] AI Config seeded")

	log.Println("  [AI] AI Knowledge Base seeding completed!")
	return nil
}
