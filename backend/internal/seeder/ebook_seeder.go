package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func (s *Seeder) seedEbooks() error {
	log.Println("  [Ebook] Seeding ebook categories & ebooks...")

	now := time.Now()

	// ══════════════════════════════════════════════════════════════
	// 1. Get existing category IDs
	// ══════════════════════════════════════════════════════════════
	catMap := map[string]uuid.UUID{}
	rows, err := s.db.Query("SELECT id, slug FROM ebook_categories WHERE is_active = true")
	if err != nil {
		return fmt.Errorf("failed to query ebook categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var slug string
		rows.Scan(&id, &slug)
		catMap[slug] = id
	}

	if len(catMap) == 0 {
		return fmt.Errorf("no ebook categories found, run migration first")
	}

	// ══════════════════════════════════════════════════════════════
	// 2. Define ebooks
	// ══════════════════════════════════════════════════════════════
	type ebookData struct {
		Title       string
		Slug        string
		Description string
		Author      string
		Price       int64
		TotalSold   int
		Categories  []string // category slugs
	}

	ebooks := []ebookData{
		{
			"Panduan Lengkap Memulai Bisnis Franchise",
			"panduan-lengkap-memulai-bisnis-franchise",
			"<p>Buku panduan komprehensif untuk Anda yang ingin memulai bisnis franchise dari nol. Membahas mulai dari pemilihan brand, analisis kelayakan, perhitungan modal, hingga strategi grand opening.</p><p>Cocok untuk pemula yang belum memiliki pengalaman bisnis sekalipun. Dilengkapi studi kasus nyata dari pengusaha franchise sukses di Indonesia.</p><h3>Apa yang akan Anda pelajari:</h3><ul><li>Cara memilih franchise yang tepat</li><li>Analisis kelayakan bisnis</li><li>Perhitungan ROI dan BEP</li><li>Strategi negosiasi dengan franchisor</li><li>Tips grand opening yang sukses</li></ul>",
			"Dr. Hendra Wijaya, MBA", 49000, 234,
			[]string{"bisnis"},
		},
		{
			"Strategi Marketing Digital untuk UMKM",
			"strategi-marketing-digital-untuk-umkm",
			"<p>Panduan praktis marketing digital khusus untuk UMKM dan bisnis kecil. Tidak perlu budget besar untuk memasarkan produk Anda secara online.</p><p>Buku ini membahas strategi yang sudah terbukti efektif: dari social media marketing, Google Ads, hingga content marketing yang bisa langsung Anda praktikkan.</p>",
			"Rina Kartika, S.Kom", 39000, 567,
			[]string{"marketing"},
		},
		{
			"Manajemen Keuangan Bisnis Outlet",
			"manajemen-keuangan-bisnis-outlet",
			"<p>Pelajari cara mengelola keuangan bisnis outlet secara profesional. Dari pencatatan sederhana hingga laporan keuangan yang bankable.</p><p>Buku ini dilengkapi template Excel siap pakai untuk pembukuan harian, mingguan, dan bulanan. Cocok untuk pemilik outlet yang ingin keuangannya lebih teratur.</p>",
			"Agus Prasetyo, SE, Ak", 55000, 189,
			[]string{"keuangan"},
		},
		{
			"SOP Operasional Outlet dari A sampai Z",
			"sop-operasional-outlet-dari-a-sampai-z",
			"<p>Standar Operasional Prosedur (SOP) lengkap untuk menjalankan outlet. Mulai dari pembukaan toko, pelayanan pelanggan, manajemen stok, hingga penutupan toko.</p><p>Dilengkapi checklist harian, mingguan, dan bulanan yang bisa langsung diadopsi untuk outlet Anda.</p>",
			"Tim BukaOutlet", 35000, 412,
			[]string{"operasional"},
		},
		{
			"Mindset Pengusaha Sukses: 7 Kebiasaan yang Mengubah Segalanya",
			"mindset-pengusaha-sukses",
			"<p>Buku motivasi bisnis yang membahas 7 kebiasaan kunci yang dimiliki oleh pengusaha sukses di Indonesia. Berdasarkan riset dan wawancara dengan 50+ entrepreneur.</p><p>Setiap bab dilengkapi action plan yang bisa langsung Anda praktikkan dalam kehidupan sehari-hari.</p>",
			"Bambang Sutrisno", 29000, 890,
			[]string{"motivasi"},
		},
		{
			"Instagram Marketing Mastery untuk Bisnis Kuliner",
			"instagram-marketing-mastery-bisnis-kuliner",
			"<p>Strategi khusus memasarkan bisnis kuliner di Instagram. Dari food photography, caption yang menjual, hashtag strategy, hingga Instagram Ads yang convert.</p><p>Dilengkapi contoh konten plan 30 hari yang bisa langsung Anda gunakan.</p>",
			"Dina Safitri", 45000, 345,
			[]string{"marketing"},
		},
		{
			"Rahasia Lokasi Strategis: Cara Memilih Tempat Usaha yang Menguntungkan",
			"rahasia-lokasi-strategis",
			"<p>Lokasi adalah faktor krusial kesuksesan outlet. Buku ini membahas framework analisis lokasi yang digunakan oleh brand-brand besar untuk memilih lokasi outlet mereka.</p><p>Pelajari cara membaca traffic, analisis kompetitor, kalkulasi sewa vs potensi omzet, dan negosiasi kontrak sewa.</p>",
			"Irfan Hakim, MT", 42000, 178,
			[]string{"bisnis", "operasional"},
		},
		{
			"Cara Menghitung BEP dan ROI Bisnis Franchise",
			"cara-menghitung-bep-dan-roi",
			"<p>Panduan praktis perhitungan Break Even Point (BEP) dan Return on Investment (ROI) khusus untuk bisnis franchise dan outlet.</p><p>Dilengkapi rumus, contoh kasus nyata, dan kalkulator Excel yang bisa langsung digunakan.</p>",
			"Prof. Sugiyono, SE, MM", 38000, 267,
			[]string{"keuangan", "bisnis"},
		},
		{
			"Food Safety & Hygiene untuk Pemilik Outlet Kuliner",
			"food-safety-hygiene-outlet-kuliner",
			"<p>Panduan standar keamanan pangan dan kebersihan yang wajib dipenuhi oleh setiap outlet kuliner. Sesuai dengan regulasi BPOM dan standar internasional HACCP.</p><p>Menghindari masalah hukum dan menjaga reputasi brand Anda.</p>",
			"Dr. Lina Marlina, M.Kes", 32000, 156,
			[]string{"operasional"},
		},
		{
			"Cara Rekrut dan Kelola Karyawan Outlet",
			"cara-rekrut-kelola-karyawan-outlet",
			"<p>Manajemen SDM praktis untuk pemilik outlet. Dari cara rekrut karyawan yang tepat, onboarding, training, evaluasi kinerja, hingga retention strategy.</p><p>Cocok untuk outlet dengan 3-20 karyawan.</p>",
			"Dewi Anggraini, S.Psi", 35000, 201,
			[]string{"operasional", "bisnis"},
		},
		{
			"TikTok Marketing untuk Bisnis Lokal",
			"tiktok-marketing-bisnis-lokal",
			"<p>Panduan lengkap memanfaatkan TikTok untuk mempromosikan bisnis lokal Anda. Dari membuat konten viral, memahami algoritma, hingga menjalankan TikTok Ads.</p><p>Studi kasus bisnis lokal yang berhasil viral dan meningkatkan omzet hingga 300%.</p>",
			"Kevin Anggara", 42000, 678,
			[]string{"marketing"},
		},
		{
			"Bisnis Franchise dengan Modal di Bawah 50 Juta",
			"bisnis-franchise-modal-dibawah-50-juta",
			"<p>Kumpulan peluang franchise dan kemitraan yang bisa dimulai dengan modal di bawah 50 juta. Dilengkapi analisis mendalam, perbandingan antar brand, dan tips memaksimalkan keuntungan.</p><p>Ideal untuk karyawan yang ingin mulai berbisnis sampingan.</p>",
			"Reza Firmansyah", 25000, 1230,
			[]string{"bisnis"},
		},
		{
			"Financial Planning untuk Entrepreneur Pemula",
			"financial-planning-entrepreneur-pemula",
			"<p>Belajar mengatur keuangan pribadi dan bisnis secara bersamaan. Banyak pengusaha pemula gagal karena mencampur keuangan pribadi dengan bisnis.</p><p>Buku ini mengajarkan cara memisahkan, mengelola, dan mengembangkan kedua sisi keuangan Anda.</p>",
			"Andri Setiawan, CFP", 48000, 345,
			[]string{"keuangan"},
		},
		{
			"Bangkit dari Kegagalan Bisnis: Kisah Nyata 10 Pengusaha",
			"bangkit-dari-kegagalan-bisnis",
			"<p>Kumpulan kisah inspiratif 10 pengusaha Indonesia yang pernah gagal total namun berhasil bangkit dan membangun bisnis yang lebih besar.</p><p>Buku ini mengajarkan bahwa kegagalan adalah bagian dari perjalanan menuju kesuksesan.</p>",
			"Yusuf Mansur", 27000, 1456,
			[]string{"motivasi"},
		},
		{
			"Strategi Pricing yang Tepat untuk Bisnis Kuliner",
			"strategi-pricing-bisnis-kuliner",
			"<p>Cara menentukan harga jual yang tepat untuk produk kuliner. Membahas cost-plus pricing, value-based pricing, psychological pricing, dan dynamic pricing.</p><p>Dilengkapi studi kasus dari berbagai jenis outlet kuliner.</p>",
			"Budi Hartono, MBA", 36000, 289,
			[]string{"keuangan", "bisnis"},
		},
		{
			"Panduan Lengkap Google My Business untuk Outlet",
			"panduan-google-my-business-outlet",
			"<p>Optimalkan kehadiran online outlet Anda dengan Google My Business. Dari setup awal, optimasi profil, manajemen review, hingga memanfaatkan fitur-fitur terbaru.</p><p>Meningkatkan visibilitas outlet Anda di Google Search dan Google Maps.</p>",
			"Sarah Amalia, S.Kom", 28000, 234,
			[]string{"marketing"},
		},
		{
			"Supply Chain Management untuk Bisnis F&B",
			"supply-chain-management-bisnis-fnb",
			"<p>Kelola rantai pasok bisnis F&B Anda secara efisien. Dari pemilihan supplier, manajemen inventory, quality control, hingga logistik distribusi.</p><p>Mengurangi waste dan meningkatkan profit margin bisnis Anda.</p>",
			"Hendarman, ST, MM", 52000, 123,
			[]string{"operasional"},
		},
		{
			"Gratis: Checklist Persiapan Buka Outlet",
			"checklist-persiapan-buka-outlet",
			"<p>Checklist komprehensif yang harus Anda persiapkan sebelum membuka outlet. Mencakup perizinan, renovasi, peralatan, karyawan, branding, dan marketing.</p><p>Ebook gratis dari Tim BukaOutlet untuk membantu Anda mempersiapkan pembukaan outlet dengan matang.</p>",
			"Tim BukaOutlet", 0, 2340,
			[]string{"bisnis", "operasional"},
		},
		{
			"Mengelola Multi-Outlet: Dari 1 Menjadi 10 Cabang",
			"mengelola-multi-outlet",
			"<p>Panduan ekspansi bisnis outlet dari single-store menjadi multi-store. Membahas standardisasi, delegasi, sistem kontrol, dan kapan waktu yang tepat untuk buka cabang baru.</p><p>Berdasarkan pengalaman nyata pemilik 10+ cabang outlet di Indonesia.</p>",
			"Herman Suharto", 58000, 167,
			[]string{"bisnis", "operasional"},
		},
		{
			"Gratis: Template Proposal Kemitraan Outlet",
			"template-proposal-kemitraan-outlet",
			"<p>Template proposal kemitraan yang profesional dan siap pakai. Cukup edit sesuai kebutuhan outlet Anda. Format yang sudah teruji dan digunakan oleh puluhan brand franchise di Indonesia.</p>",
			"Tim BukaOutlet", 0, 3120,
			[]string{"bisnis"},
		},
	}

	// ══════════════════════════════════════════════════════════════
	// 3. Insert ebooks
	// ══════════════════════════════════════════════════════════════
	ebookCount := 0
	for _, eb := range ebooks {
		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM ebooks WHERE slug = $1)", eb.Slug).Scan(&exists)
		if exists {
			continue
		}

		ebookID := uuid.New()

		_, err := s.db.Exec(`
			INSERT INTO ebooks (id, title, slug, description, author, price, is_active, total_sold, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, true, $7, $8, $8)`,
			ebookID, eb.Title, eb.Slug, eb.Description, eb.Author, eb.Price, eb.TotalSold, now,
		)
		if err != nil {
			return fmt.Errorf("failed to seed ebook %s: %w", eb.Title, err)
		}

		// Insert category mappings
		for _, catSlug := range eb.Categories {
			catID, ok := catMap[catSlug]
			if !ok {
				continue
			}
			_, err := s.db.Exec(
				"INSERT INTO ebook_category_mapping (ebook_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				ebookID, catID,
			)
			if err != nil {
				return fmt.Errorf("failed to map ebook category %s->%s: %w", eb.Slug, catSlug, err)
			}
		}

		ebookCount++
	}

	log.Printf("    [OK] %d ebooks seeded", ebookCount)
	log.Println("  [Ebook] Ebook seeding completed!")
	return nil
}
