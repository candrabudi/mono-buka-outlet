package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (s *Seeder) seedOutlets() error {
	log.Println("  [Outlet] Seeding outlet categories, outlets & packages...")

	now := time.Now()

	// ══════════════════════════════════════════════════════════════
	// 1. Get existing category IDs
	// ══════════════════════════════════════════════════════════════
	catMap := map[string]uuid.UUID{}
	rows, err := s.db.Query("SELECT id, slug FROM outlet_categories WHERE is_active = true")
	if err != nil {
		return fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var slug string
		rows.Scan(&id, &slug)
		catMap[slug] = id
	}

	if len(catMap) == 0 {
		return fmt.Errorf("no outlet categories found, run migration first")
	}

	// ══════════════════════════════════════════════════════════════
	// 2. Define 30 outlets
	// ══════════════════════════════════════════════════════════════
	type outletData struct {
		Name                string
		Slug                string
		CategorySlug        string
		ShortDescription    string
		Description         string
		MinInvestment       float64
		MaxInvestment       float64
		ProfitSharing       float64
		EstimatedROI        string
		LocationRequirement string
		City                string
		Province            string
		ContactWhatsapp     string
		ContactEmail        string
		IsFeatured          bool
		TotalOutlets        int
		YearEstablished     int
	}

	outlets := []outletData{
		{"Ayam Geprek Nusantara", "ayam-geprek-nusantara", "franchise", "Franchise ayam geprek terlaris se-Indonesia", "Ayam Geprek Nusantara adalah franchise kuliner ayam geprek dengan resep rahasia bumbu rempah khas Indonesia. Sudah terbukti dengan 150+ outlet di seluruh Indonesia, omzet rata-rata 80-120 juta per bulan.", 75000000, 150000000, 30, "8-12 bulan", "Lokasi strategis di area komersial, minimal 40m2", "Jakarta Selatan", "DKI Jakarta", "08123456001", "info@ayamgepreknusantara.id", true, 156, 2018},
		{"Kopi Kenangan Lokal", "kopi-kenangan-lokal", "franchise", "Coffee shop kekinian dengan harga terjangkau", "Kopi Kenangan Lokal hadir sebagai solusi bisnis kopi grab-and-go dengan menu kopi susu kekinian, teh, dan snack. Konsep minimalis modern cocok untuk area perkantoran dan kampus.", 120000000, 250000000, 25, "10-14 bulan", "Area perkantoran/kampus, minimal 20m2", "Jakarta Pusat", "DKI Jakarta", "08123456002", "franchise@kopikenanganlokal.id", true, 89, 2019},
		{"Martabak Mertua", "martabak-mertua", "kemitraan", "Kemitraan martabak premium rasa rumahan", "Martabak Mertua menyajikan martabak manis dan telur dengan resep turun-temurun. Bahan premium, porsi jumbo, harga bersahabat. Cocok untuk area residensial dan kuliner malam.", 35000000, 80000000, 35, "6-10 bulan", "Lokasi ramai, area kuliner/residensial, minimal 15m2", "Bandung", "Jawa Barat", "08123456003", "partner@martabakmertua.id", true, 210, 2016},
		{"Bakso Boedjang", "bakso-boedjang", "franchise", "Franchise bakso urat legendaris Bandung", "Bakso Boedjang terkenal dengan bakso urat jumbo dan kuah kaldu sapi premium. Sudah hadir di 5 provinsi dengan antrian panjang setiap hari.", 90000000, 200000000, 28, "9-13 bulan", "Ruko atau stand alone, minimal 60m2, parkir memadai", "Bandung", "Jawa Barat", "08123456004", "info@baksoboedjang.id", false, 67, 2017},
		{"Sate Taichan Goreng", "sate-taichan-goreng", "kemitraan", "Sate taichan viral dengan sambal khas", "Sate Taichan Goreng menawarkan konsep sate ayam tanpa bumbu kacang, diganti sambal matah dan sambal hijau khas. Operasional simpel, cocok untuk pemula.", 25000000, 50000000, 40, "4-7 bulan", "Stand/booth di area kuliner, minimal 6m2", "Surabaya", "Jawa Timur", "08123456005", "info@satetaichan.id", true, 340, 2020},
		{"Dimsum Go!", "dimsum-go", "franchise", "Franchise dimsum frozen dan fresh", "Dimsum Go! menyediakan berbagai jenis dimsum fresh dan frozen dengan kualitas restoran. Cocok untuk model grab-and-go atau dine-in kecil.", 60000000, 120000000, 30, "8-11 bulan", "Mall, food court, atau ruko, minimal 25m2", "Jakarta Barat", "DKI Jakarta", "08123456006", "franchise@dimsumgo.id", false, 45, 2021},
		{"Nasi Goreng Sultan", "nasi-goreng-sultan", "kemitraan", "Nasi goreng premium dengan topping mewah", "Nasi Goreng Sultan hadir dengan konsep nasi goreng premium topping wagyu, seafood, dan cheese. Viral di media sosial, cocok untuk area nightlife dan residensial.", 20000000, 45000000, 45, "3-6 bulan", "Gerobak atau booth, minimal 6m2", "Yogyakarta", "DI Yogyakarta", "08123456007", "info@nasgorengsultan.id", true, 180, 2021},
		{"Es Teh Solo", "es-teh-solo", "franchise", "Franchise minuman teh khas Solo", "Es Teh Solo berfokus pada minuman teh dengan berbagai varian: teh tarik, teh susu, thai tea, lemon tea. Modal kecil, operasi simpel, margin tinggi.", 30000000, 60000000, 35, "4-8 bulan", "Booth/container, minimal 4m2, area ramai", "Solo", "Jawa Tengah", "08123456008", "franchise@estehsolo.id", true, 520, 2019},
		{"Pizza Lokal", "pizza-lokal", "licensing", "Lisensi brand pizza dengan cita rasa lokal", "Pizza Lokal menggabungkan teknik pizza Italia dengan bumbu dan topping khas Indonesia seperti rendang, sate, dan sambal matah. Konsep unik yang menarik perhatian.", 150000000, 300000000, 22, "12-18 bulan", "Ruko atau standalone, minimal 80m2, dapur lengkap", "Jakarta Selatan", "DKI Jakarta", "08123456009", "license@pizzalokal.id", false, 23, 2020},
		{"Jamu Modern Nyonya", "jamu-modern-nyonya", "kemitraan", "Jamu tradisional dalam kemasan modern", "Jamu Modern Nyonya meramu jamu tradisional Indonesia dalam kemasan modern dan instagramable. Varian: beras kencur, kunyit asam, temulawak latte, jahe merah.", 15000000, 30000000, 50, "3-5 bulan", "Booth kecil, minimal 3m2, area perkantoran", "Semarang", "Jawa Tengah", "08123456010", "info@jamunyonya.id", true, 150, 2022},
		{"Burger Blenger", "burger-blenger", "franchise", "Franchise burger smashed viral", "Burger Blenger terkenal dengan smashed burger juicy dengan saus secret recipe. Menu simple tapi addictive. Sudah terbukti viral di sosmed dengan jutaan views.", 80000000, 160000000, 30, "8-12 bulan", "Ruko atau booth besar, minimal 30m2", "Jakarta Timur", "DKI Jakarta", "08123456011", "franchise@burgerblenger.id", true, 78, 2021},
		{"Roti Bakar Bandung", "roti-bakar-bandung", "kemitraan", "Roti bakar khas Bandung dengan beragam topping", "Roti Bakar Bandung menyajikan roti bakar premium dengan puluhan pilihan topping: coklat, keju, green tea, taro, ovomaltine, dan lainnya.", 18000000, 40000000, 40, "4-6 bulan", "Gerobak atau booth, area kuliner malam", "Bandung", "Jawa Barat", "08123456012", "info@rotibakarbdg.id", false, 290, 2017},
		{"Mie Gacoan", "mie-gacoan", "franchise", "Franchise mie pedas level terlaris", "Mie Gacoan hadir dengan mie pedas berlevel yang sangat populer di kalangan anak muda. Harga terjangkau, porsi besar, dan tempat instagramable.", 200000000, 500000000, 20, "14-20 bulan", "Ruko besar/standalone, minimal 150m2, parkir luas", "Malang", "Jawa Timur", "08123456013", "franchise@miegacoan.id", true, 120, 2016},
		{"Kebab Turki Baba", "kebab-turki-baba", "franchise", "Franchise kebab ukuran jumbo", "Kebab Turki Baba menyajikan kebab autentik Turki dengan daging premium, sayuran segar, dan saus khas. Ukuran jumbo menjadi ciri khas yang membedakan.", 45000000, 90000000, 32, "6-10 bulan", "Booth/gerobak, minimal 6m2, area ramai", "Surabaya", "Jawa Timur", "08123456014", "info@kebabtbaba.id", false, 180, 2015},
		{"Chicken Crispy Bang", "chicken-crispy-bang", "kemitraan", "Ayam crispy ala Korea yang terjangkau", "Chicken Crispy Bang menyajikan ayam goreng crispy ala Korean fried chicken dengan saus gochujang, honey butter, dan cheese. Harga student-friendly.", 22000000, 55000000, 38, "5-8 bulan", "Booth atau ruko kecil, minimal 15m2", "Depok", "Jawa Barat", "08123456015", "info@chickencrispybang.id", true, 95, 2022},
		{"Warung Steak Murah", "warung-steak-murah", "franchise", "Steak premium harga warung", "Warung Steak Murah mematahkan stigma steak mahal. Dengan daging impor berkualitas dan teknik masak profesional, steak dijual mulai 25 ribu.", 100000000, 220000000, 25, "10-15 bulan", "Ruko, minimal 80m2, area komersial", "Medan", "Sumatera Utara", "08123456016", "franchise@warungsteak.id", false, 55, 2018},
		{"Takoyaki Corner", "takoyaki-corner", "kemitraan", "Takoyaki Jepang dengan isian lokal", "Takoyaki Corner menyajikan takoyaki Jepang dengan twist lokal: isian bakso, sosis, keju mozarella, hingga rendang. Operasi simpel dengan 1-2 orang.", 12000000, 25000000, 45, "2-4 bulan", "Stand kecil, minimal 3m2", "Bekasi", "Jawa Barat", "08123456017", "info@takoyakicorner.id", true, 410, 2020},
		{"Salad Bar Fresh", "salad-bar-fresh", "licensing", "Lisensi salad bar sehat dan segar", "Salad Bar Fresh menyediakan salad custom dengan puluhan pilihan topping dan dressing. Target market: kantoran, gym-goers, dan health-conscious community.", 80000000, 150000000, 28, "10-14 bulan", "Mall atau area perkantoran, minimal 25m2", "Jakarta Selatan", "DKI Jakarta", "08123456018", "license@saladbarfresh.id", false, 18, 2021},
		{"Bubur Ayam Cirebon", "bubur-ayam-cirebon", "kemitraan", "Bubur ayam khas Cirebon dengan cakwe homemade", "Bubur Ayam Cirebon merupakan kemitraan bubur ayam dengan resep autentik Cirebon. Cakwe homemade, ayam suwir berbumbu, dan kuah kaldu premium.", 20000000, 45000000, 38, "4-7 bulan", "Gerobak atau ruko kecil, minimal 10m2", "Cirebon", "Jawa Barat", "08123456019", "info@buburayamcirebon.id", true, 130, 2019},
		{"Boba Time!", "boba-time", "franchise", "Franchise boba drink premium", "Boba Time! menghadirkan minuman boba premium dengan boba homemade yang kenyal. Varian: brown sugar, taro, matcha, strawberry, dan seasonal specials.", 55000000, 110000000, 30, "7-10 bulan", "Mall, area kampus/sekolah, minimal 8m2", "Tangerang", "Banten", "08123456020", "franchise@bobatime.id", true, 200, 2020},
		{"Seafood Kiloan Pak De", "seafood-kiloan-pak-de", "franchise", "Franchise seafood segar sistem kiloan", "Seafood Kiloan Pak De menyajikan seafood segar yang dipilih langsung oleh pelanggan dengan sistem kiloan. Bumbu racik khas Jawa dengan sambal pilihan.", 130000000, 280000000, 25, "12-16 bulan", "Ruko/bangunan, minimal 100m2, area strategis", "Semarang", "Jawa Tengah", "08123456021", "franchise@seafoodpakde.id", false, 35, 2019},
		{"Donat Kampung", "donat-kampung", "kemitraan", "Donat lembut harga kampung rasa premium", "Donat Kampung memproduksi donat lembut dengan harga sangat terjangkau (mulai 3 ribu). Resep premium, bahan berkualitas, tapi tetap merakyat.", 10000000, 22000000, 50, "2-4 bulan", "Booth kecil atau titip jual, minimal 2m2", "Bogor", "Jawa Barat", "08123456022", "info@donatkampung.id", true, 680, 2018},
		{"Tahu Crispy Lestari", "tahu-crispy-lestari", "kemitraan", "Tahu crispy renyah dengan saus kacang spesial", "Tahu Crispy Lestari menjual tahu crispy dengan bumbu rempah dan saus kacang spesial. Modal sangat kecil, margin besar, operasional mudah.", 8000000, 15000000, 55, "1-3 bulan", "Gerobak atau etalase kecil, minimal 2m2", "Tasikmalaya", "Jawa Barat", "08123456023", "info@tahulestari.id", false, 450, 2017},
		{"Rice Bowl Kita", "rice-bowl-kita", "franchise", "Franchise rice bowl ala restoran harga mahasiswa", "Rice Bowl Kita menyajikan rice bowl dengan lauk premium: chicken teriyaki, beef bulgogi, crispy chicken, salmon. Sajian ala restoran, harga mahasiswa.", 50000000, 100000000, 32, "6-9 bulan", "Food court, mall, area kampus, minimal 15m2", "Yogyakarta", "DI Yogyakarta", "08123456024", "franchise@ricebowlkita.id", true, 110, 2021},
		{"Warung Nasi Padang Mini", "warung-nasi-padang-mini", "bot", "BOT warung nasi Padang skala mini", "Warung Nasi Padang Mini menghadirkan konsep nasi Padang dalam format yang lebih compact dan modern. Menu andalan: rendang, ayam bakar, gulai nangka.", 70000000, 140000000, 30, "8-12 bulan", "Ruko kecil, minimal 30m2, area perkantoran", "Padang", "Sumatera Barat", "08123456025", "info@naspadangmini.id", true, 42, 2020},
		{"Thai Tea Mama", "thai-tea-mama", "kemitraan", "Thai tea original dengan resep asli Thailand", "Thai Tea Mama menggunakan daun teh impor Thailand original dengan susu segar. Rasa authentic yang sulit ditiru kompetitor. Modal kecil, untung besar.", 12000000, 28000000, 45, "3-5 bulan", "Booth kecil, minimal 3m2", "Makassar", "Sulawesi Selatan", "08123456026", "info@thaiteamama.id", true, 320, 2021},
		{"Sushi Roll Express", "sushi-roll-express", "licensing", "Lisensi sushi roll grab-and-go", "Sushi Roll Express membawa konsep sushi affordable dalam format grab-and-go. Menu utama: salmon roll, tuna roll, ebi furai roll, chicken katsu roll.", 100000000, 200000000, 25, "10-15 bulan", "Mall atau area premium, minimal 20m2", "Jakarta Pusat", "DKI Jakarta", "08123456027", "license@sushirollexpress.id", false, 15, 2022},
		{"Ayam Bakar Wong Solo", "ayam-bakar-wong-solo", "franchise", "Franchise ayam bakar bumbu rempah legendaris", "Ayam Bakar Wong Solo merupakan franchise ayam bakar dengan bumbu rempah khas Solo yang sudah legendaris sejak 2005. Menu lengkap: ayam bakar, nasi liwet, pecel.", 180000000, 400000000, 22, "14-20 bulan", "Ruko/bangunan besar, minimal 120m2, parkir luas", "Solo", "Jawa Tengah", "08123456028", "franchise@wongsolo.id", true, 85, 2005},
		{"Crepes & Pancake House", "crepes-pancake-house", "franchise", "Franchise crepes dan pancake premium", "Crepes & Pancake House menyajikan crepes tipis renyah dan fluffy pancake ala Jepang. Pilihan topping manis dan savory. Konsep cafe kecil yang instagramable.", 65000000, 130000000, 30, "7-11 bulan", "Mall atau area lifestyle, minimal 25m2", "Surabaya", "Jawa Timur", "08123456029", "franchise@crepeshouse.id", false, 30, 2022},
		{"Pecel Lele Lamongan", "pecel-lele-lamongan", "kemitraan", "Pecel lele autentik Lamongan dengan sambal terasi", "Pecel Lele Lamongan hadir dengan rasa otentik khas Lamongan. Lele goreng crispy, sambal terasi bawang, lalapan segar, dan nasi hangat. Harga terjangkau untuk semua kalangan.", 15000000, 35000000, 42, "3-6 bulan", "Tenda atau ruko kecil, minimal 10m2, area padat penduduk", "Lamongan", "Jawa Timur", "08123456030", "info@pecellelelamongan.id", true, 380, 2016},
	}

	// ══════════════════════════════════════════════════════════════
	// 3. Insert outlets
	// ══════════════════════════════════════════════════════════════
	outletIDs := make([]uuid.UUID, len(outlets))

	for i, o := range outlets {
		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM outlets WHERE slug = $1 AND deleted_at IS NULL)", o.Slug).Scan(&exists)
		if exists {
			// Get existing ID for package seeding
			_ = s.db.QueryRow("SELECT id FROM outlets WHERE slug = $1 AND deleted_at IS NULL", o.Slug).Scan(&outletIDs[i])
			continue
		}

		outletID := uuid.New()
		outletIDs[i] = outletID

		catID := catMap[o.CategorySlug]

		_, err := s.db.Exec(`
			INSERT INTO outlets (
				id, name, slug, category, category_id, short_description, description,
				minimum_investment, maximum_investment, profit_sharing_percentage,
				estimated_roi, location_requirement, city, province,
				contact_whatsapp, contact_email, is_active, is_featured,
				total_outlets, year_established, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10, $11, $12, $13, $14,
				$15, $16, $17, $18, $19, $20, $21, $21
			)`,
			outletID, o.Name, o.Slug, o.CategorySlug, catID, o.ShortDescription, o.Description,
			o.MinInvestment, o.MaxInvestment, o.ProfitSharing,
			o.EstimatedROI, o.LocationRequirement, o.City, o.Province,
			o.ContactWhatsapp, o.ContactEmail, true, o.IsFeatured,
			o.TotalOutlets, o.YearEstablished, now,
		)
		if err != nil {
			return fmt.Errorf("failed to seed outlet %s: %w", o.Name, err)
		}
	}

	log.Printf("    [OK] %d outlets seeded", len(outlets))

	// ══════════════════════════════════════════════════════════════
	// 4. Define packages per outlet
	// ══════════════════════════════════════════════════════════════
	type pkgData struct {
		OutletIdx    int
		Name         string
		Slug         string
		Price        int64
		MinimumDP    int64
		Duration     string
		EstimatedBEP string
		NetProfit    string
		Description  string
		Benefits     []string
		SortOrder    int
	}

	packages := []pkgData{
		// Ayam Geprek Nusantara (0)
		{0, "Paket Booth", "paket-booth", 75000000, 30000000, "5 tahun", "8-10 bulan", "15-25 juta/bulan", "Paket booth untuk area food court atau mall", []string{"Booth portable", "Peralatan masak lengkap", "Bahan baku awal", "Training 7 hari", "SOP operasional"}, 1},
		{0, "Paket Ruko", "paket-ruko", 150000000, 50000000, "5 tahun", "10-14 bulan", "25-40 juta/bulan", "Paket ruko untuk area komersial", []string{"Desain interior", "Peralatan masak lengkap", "Bahan baku awal 1 bulan", "Training 14 hari", "SOP operasional", "Support marketing 3 bulan"}, 2},

		// Kopi Kenangan Lokal (1)
		{1, "Paket Grab & Go", "paket-grab-go", 120000000, 40000000, "5 tahun", "10-12 bulan", "20-30 juta/bulan", "Konsep take-away modern compact", []string{"Mesin espresso premium", "Grinder profesional", "Bahan baku 2 minggu", "Training barista 10 hari", "Desain booth"}, 1},
		{1, "Paket Cafe", "paket-cafe", 250000000, 80000000, "5 tahun", "12-16 bulan", "35-55 juta/bulan", "Konsep cafe dengan dine-in area", []string{"Mesin espresso premium", "Full interior design", "Furniture set", "Training barista 14 hari", "Marketing support 6 bulan", "Bahan baku 1 bulan"}, 2},

		// Martabak Mertua (2)
		{2, "Paket Gerobak", "paket-gerobak", 35000000, 15000000, "3 tahun", "5-7 bulan", "8-15 juta/bulan", "Paket gerobak martabak lengkap", []string{"Gerobak custom", "Peralatan masak", "Bahan baku awal", "Training 5 hari", "Resep lengkap"}, 1},
		{2, "Paket Premium", "paket-premium", 80000000, 30000000, "5 tahun", "8-12 bulan", "15-28 juta/bulan", "Paket booth premium dengan etalase", []string{"Booth premium", "Peralatan masak lengkap", "Etalase display", "Bahan baku 2 minggu", "Training 7 hari", "Branding material"}, 2},

		// Bakso Boedjang (3)
		{3, "Paket Warung", "paket-warung", 90000000, 35000000, "5 tahun", "8-11 bulan", "18-30 juta/bulan", "Paket warung bakso sederhana", []string{"Peralatan masak", "Meja kursi 20 set", "Bahan baku 2 minggu", "Training 10 hari", "SOP resep"}, 1},
		{3, "Paket Resto", "paket-resto", 200000000, 70000000, "7 tahun", "12-15 bulan", "35-50 juta/bulan", "Paket restoran bakso premium", []string{"Interior design", "Peralatan lengkap", "Furniture 40 set", "Training 14 hari", "Marketing support", "AC & peralatan"}, 2},

		// Sate Taichan Goreng (4)
		{4, "Paket Starter", "paket-starter", 25000000, 10000000, "3 tahun", "3-5 bulan", "5-10 juta/bulan", "Mulai usaha sate taichan dengan modal kecil", []string{"Gerobak/booth", "Peralatan grill", "Bahan baku awal", "Training 3 hari"}, 1},
		{4, "Paket Complete", "paket-complete", 50000000, 20000000, "5 tahun", "5-8 bulan", "10-20 juta/bulan", "Paket lengkap dengan booth premium", []string{"Booth premium", "Peralatan grill premium", "Bahan baku 2 minggu", "Training 5 hari", "Branding material", "Marketing support"}, 2},

		// Dimsum Go! (5)
		{5, "Paket Express", "paket-express", 60000000, 25000000, "5 tahun", "7-9 bulan", "12-20 juta/bulan", "Konsep grab-and-go dimsum", []string{"Steamer set", "Display counter", "Packaging 1 bulan", "Training 7 hari", "Bahan baku awal"}, 1},
		{5, "Paket Dine-in", "paket-dinein", 120000000, 45000000, "5 tahun", "10-13 bulan", "22-35 juta/bulan", "Konsep dine-in dengan area makan", []string{"Full kitchen set", "Interior design", "Furniture", "Training 10 hari", "Bahan baku 2 minggu", "Marketing support"}, 2},

		// Nasi Goreng Sultan (6)
		{6, "Paket Sultan Mini", "paket-sultan-mini", 20000000, 8000000, "3 tahun", "2-4 bulan", "5-12 juta/bulan", "Gerobak nasi goreng premium", []string{"Gerobak custom", "Kompor & wajan besar", "Bahan baku awal", "Resep rahasia", "Training 3 hari"}, 1},
		{6, "Paket Sultan Max", "paket-sultan-max", 45000000, 18000000, "5 tahun", "4-7 bulan", "12-22 juta/bulan", "Booth nasi goreng premium", []string{"Booth premium", "Peralatan lengkap", "Bahan baku 2 minggu", "Training 5 hari", "Branding lengkap"}, 2},

		// Es Teh Solo (7)
		{7, "Paket Mini", "paket-mini", 30000000, 12000000, "3 tahun", "3-5 bulan", "7-12 juta/bulan", "Booth teh kecil untuk area ramai", []string{"Booth portable", "Peralatan brewing", "Cup & straw 1 bulan", "Bahan baku 2 minggu", "Training 3 hari"}, 1},
		{7, "Paket Regular", "paket-regular", 60000000, 25000000, "5 tahun", "5-9 bulan", "12-22 juta/bulan", "Booth teh premium di mall/area komersial", []string{"Booth premium", "Mesin brewing", "Packaging 1 bulan", "Bahan baku 1 bulan", "Training 5 hari", "Marketing support"}, 2},

		// Pizza Lokal (8)
		{8, "Paket Standard", "paket-standard", 150000000, 50000000, "5 tahun", "12-14 bulan", "25-40 juta/bulan", "Restoran pizza lokal standard", []string{"Pizza oven", "Kitchen equipment", "Interior basic", "Training pizzaiolo 14 hari", "Bahan baku awal"}, 1},
		{8, "Paket Premium", "paket-premium", 300000000, 100000000, "7 tahun", "15-20 bulan", "45-70 juta/bulan", "Restoran pizza premium full service", []string{"Double pizza oven", "Full kitchen", "Interior design premium", "Furniture set", "Training 21 hari", "Marketing 6 bulan"}, 2},

		// Jamu Modern Nyonya (9)
		{9, "Paket Jamu Go", "paket-jamu-go", 15000000, 6000000, "3 tahun", "2-4 bulan", "4-8 juta/bulan", "Booth jamu modern kecil", []string{"Booth mini", "Blender & peralatan", "Botol & cup 500pcs", "Bahan baku 2 minggu", "Resep 10 varian"}, 1},
		{9, "Paket Nyonya", "paket-nyonya", 30000000, 12000000, "5 tahun", "4-6 bulan", "8-15 juta/bulan", "Booth jamu premium", []string{"Booth premium", "Peralatan lengkap", "Packaging premium", "Bahan baku 1 bulan", "Training 5 hari", "Resep 20 varian"}, 2},

		// Burger Blenger (10)
		{10, "Paket Smash", "paket-smash", 80000000, 30000000, "5 tahun", "7-10 bulan", "15-25 juta/bulan", "Booth burger smashed compact", []string{"Griddle premium", "Booth", "Bahan baku 2 minggu", "Training 7 hari", "Packaging 1 bulan"}, 1},
		{10, "Paket Blenger", "paket-blenger", 160000000, 55000000, "5 tahun", "10-14 bulan", "28-42 juta/bulan", "Restoran burger premium", []string{"Full kitchen", "Interior design", "Furniture set", "Training 10 hari", "Marketing support 3 bulan", "Bahan baku 1 bulan"}, 2},

		// Donat Kampung (21)
		{21, "Paket Rintisan", "paket-rintisan", 10000000, 5000000, "2 tahun", "1-2 bulan", "3-6 juta/bulan", "Mulai jualan donat kampung dari rumah", []string{"Peralatan produksi", "Bahan baku awal", "Resep 10 varian", "Training 2 hari", "Etalase display"}, 1},
		{21, "Paket Toko", "paket-toko", 22000000, 10000000, "3 tahun", "3-5 bulan", "6-12 juta/bulan", "Toko donat kampung lengkap", []string{"Peralatan produksi lengkap", "Etalase premium", "Bahan baku 2 minggu", "Training 4 hari", "Branding material", "Packaging 1 bulan"}, 2},

		// Mie Gacoan (12)
		{12, "Paket Standard", "paket-standard", 200000000, 70000000, "7 tahun", "14-16 bulan", "35-50 juta/bulan", "Restoran mie gacoan standard", []string{"Full kitchen setup", "Interior design", "Furniture 60 set", "Training 14 hari", "Bahan baku 2 minggu"}, 1},
		{12, "Paket Flagship", "paket-flagship", 500000000, 150000000, "10 tahun", "18-24 bulan", "60-100 juta/bulan", "Restoran mie gacoan flagship premium", []string{"Premium interior", "Full kitchen", "Furniture 100 set", "Training 21 hari", "Grand opening support", "Marketing 6 bulan", "Outdoor area"}, 2},

		// Boba Time! (19)
		{19, "Paket Boba Mini", "paket-boba-mini", 55000000, 20000000, "5 tahun", "6-8 bulan", "10-18 juta/bulan", "Booth boba minimalis", []string{"Sealer & shaker", "Booth premium", "Cup & straw 2000pcs", "Bahan baku 2 minggu", "Training 5 hari"}, 1},
		{19, "Paket Boba Max", "paket-boba-max", 110000000, 40000000, "5 tahun", "9-12 bulan", "20-32 juta/bulan", "Toko boba premium dengan dine-in", []string{"Full equipment", "Interior design", "Furniture", "Training 7 hari", "Marketing support", "Bahan baku 1 bulan", "Packaging 2 bulan"}, 2},

		// Rice Bowl Kita (23)
		{23, "Paket Food Court", "paket-food-court", 50000000, 20000000, "5 tahun", "5-7 bulan", "10-18 juta/bulan", "Stand rice bowl di food court", []string{"Kitchen equipment", "Counter display", "Packaging 1 bulan", "Bahan baku 2 minggu", "Training 5 hari"}, 1},
		{23, "Paket Cafe Bowl", "paket-cafe-bowl", 100000000, 35000000, "5 tahun", "8-11 bulan", "20-32 juta/bulan", "Cafe rice bowl dengan dine-in", []string{"Full kitchen", "Interior design", "Furniture", "Training 10 hari", "Marketing support", "Bahan baku 1 bulan"}, 2},

		// Warung Nasi Padang Mini (24)
		{24, "Paket Warung", "paket-warung", 70000000, 25000000, "5 tahun", "7-10 bulan", "12-22 juta/bulan", "Warung nasi Padang compact", []string{"Etalase Padang", "Peralatan masak", "Resep lengkap", "Training 10 hari", "Bahan baku awal"}, 1},
		{24, "Paket Premium", "paket-premium", 140000000, 50000000, "7 tahun", "10-14 bulan", "25-40 juta/bulan", "Restoran nasi Padang modern", []string{"Interior modern", "Etalase premium", "Full kitchen", "Training 14 hari", "Marketing support", "Bahan baku 2 minggu"}, 2},

		// Thai Tea Mama (25)
		{25, "Paket Booth", "paket-booth", 12000000, 5000000, "3 tahun", "2-3 bulan", "4-8 juta/bulan", "Booth thai tea mini", []string{"Booth", "Peralatan brewing", "Cup 500pcs", "Bahan baku 2 minggu", "Resep 8 varian"}, 1},
		{25, "Paket Premium", "paket-premium", 28000000, 12000000, "5 tahun", "4-6 bulan", "8-15 juta/bulan", "Booth thai tea premium", []string{"Booth premium", "Full equipment", "Cup & packaging 1 bulan", "Training 3 hari", "Bahan baku 1 bulan", "Branding lengkap"}, 2},

		// Pecel Lele Lamongan (29)
		{29, "Paket Tenda", "paket-tenda", 15000000, 6000000, "3 tahun", "2-4 bulan", "5-10 juta/bulan", "Tenda pecel lele malam hari", []string{"Tenda & meja kursi", "Peralatan masak", "Bahan baku awal", "Resep sambal 5 varian", "Training 3 hari"}, 1},
		{29, "Paket Warung", "paket-warung", 35000000, 15000000, "5 tahun", "4-7 bulan", "10-18 juta/bulan", "Warung pecel lele full-time", []string{"Setup warung", "Peralatan lengkap", "Meja kursi 15 set", "Bahan baku 2 minggu", "Training 5 hari", "Spanduk & branding"}, 2},
	}

	// ══════════════════════════════════════════════════════════════
	// 5. Insert packages
	// ══════════════════════════════════════════════════════════════
	pkgCount := 0
	for _, p := range packages {
		if outletIDs[p.OutletIdx] == uuid.Nil {
			continue
		}

		var exists bool
		_ = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM outlet_packages WHERE outlet_id = $1 AND slug = $2)",
			outletIDs[p.OutletIdx], p.Slug).Scan(&exists)
		if exists {
			continue
		}

		_, err := s.db.Exec(`
			INSERT INTO outlet_packages (
				id, outlet_id, name, slug, price, minimum_dp,
				duration, estimated_bep, net_profit, description,
				benefits, sort_order, is_active, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10,
				$11, $12, true, $13, $13
			)`,
			uuid.New(), outletIDs[p.OutletIdx], p.Name, p.Slug, p.Price, p.MinimumDP,
			p.Duration, p.EstimatedBEP, p.NetProfit, p.Description,
			pq.Array(p.Benefits), p.SortOrder, now,
		)
		if err != nil {
			return fmt.Errorf("failed to seed package %s for outlet idx %d: %w", p.Name, p.OutletIdx, err)
		}
		pkgCount++
	}

	log.Printf("    [OK] %d outlet packages seeded", pkgCount)
	log.Println("  [Outlet] Outlet seeding completed!")
	return nil
}
