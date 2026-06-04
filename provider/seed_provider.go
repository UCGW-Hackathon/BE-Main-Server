package provider

import (
	"log"
	"situkang/models/entity"
	"time"

	"gorm.io/gorm"
)

// seedReferenceData adalah entry point utama untuk semua seeding.
// Dipanggil dari config_provider.go saat startup.
func seedReferenceData(db *gorm.DB) error {
	log.Println("[Seed] Seeding categories...")
	if err := seedCategories(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding base services...")
	if err := seedBaseServices(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding extended services...")
	if err := SeedExtendedServices(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding articles...")
	if err := seedArticles(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding FAQs...")
	if err := seedFAQs(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding promotions...")
	if err := seedPromotions(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding workers & customers...")
	if err := SeedWorkers(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding orders, chat & history...")
	if err := SeedOrdersAndHistory(db); err != nil {
		return err
	}
	log.Println("[Seed] Seeding notifications...")
	if err := seedNotifications(db); err != nil {
		return err
	}
	log.Println("[Seed] ✅ All seeds completed successfully.")
	return nil
}

// ─── CATEGORIES ───────────────────────────────────────────────────────────────

func seedCategories(db *gorm.DB) error {
	categories := []entity.Category{
		{Name: "AC", Slug: "ac", IconURL: strPtr("https://cdn.handydirect.id/icons/ac.png"), Description: strPtr("Servis, pemasangan, dan perbaikan AC"), DisplayOrder: 1, IsActive: true},
		{Name: "Pipa", Slug: "pipa", IconURL: strPtr("https://cdn.handydirect.id/icons/pipa.png"), Description: strPtr("Perbaikan dan pemasangan pipa air"), DisplayOrder: 2, IsActive: true},
		{Name: "Atap", Slug: "atap", IconURL: strPtr("https://cdn.handydirect.id/icons/atap.png"), Description: strPtr("Perbaikan atap bocor dan rangka atap"), DisplayOrder: 3, IsActive: true},
		{Name: "Listrik", Slug: "listrik", IconURL: strPtr("https://cdn.handydirect.id/icons/listrik.png"), Description: strPtr("Instalasi dan perbaikan kelistrikan rumah"), DisplayOrder: 4, IsActive: true},
		{Name: "Kunci", Slug: "kunci", IconURL: strPtr("https://cdn.handydirect.id/icons/kunci.png"), Description: strPtr("Perbaikan kunci, pintu, dan keamanan rumah"), DisplayOrder: 5, IsActive: true},
		{Name: "Kayu", Slug: "kayu", IconURL: strPtr("https://cdn.handydirect.id/icons/kayu.png"), Description: strPtr("Pekerjaan kayu, furniture, dan perbaikan kusen"), DisplayOrder: 6, IsActive: true},
		{Name: "Cat", Slug: "cat", IconURL: strPtr("https://cdn.handydirect.id/icons/cat.png"), Description: strPtr("Pengecatan interior dan eksterior"), DisplayOrder: 7, IsActive: true},
		{Name: "Kebun", Slug: "kebun", IconURL: strPtr("https://cdn.handydirect.id/icons/kebun.png"), Description: strPtr("Perawatan taman dan kebun rumah"), DisplayOrder: 8, IsActive: true},
	}
	for _, c := range categories {
		var existing entity.Category
		if err := db.Where("slug = ?", c.Slug).Assign(c).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── BASE SERVICES ────────────────────────────────────────────────────────────

func seedBaseServices(db *gorm.DB) error {
	type svcSeed struct {
		CatSlug string
		Svc     entity.Service
	}
	seeds := []svcSeed{
		{"ac", entity.Service{Name: "Servis AC", Slug: "servis-ac", Description: strPtr("Cuci, tambah freon, dan perbaikan AC"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
		{"pipa", entity.Service{Name: "Perbaikan Pipa", Slug: "perbaikan-pipa", Description: strPtr("Perbaikan pipa bocor dan saluran air"), BasePrice: intPtr(120000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},
		{"atap", entity.Service{Name: "Perbaikan Atap Bocor", Slug: "perbaikan-atap-bocor", Description: strPtr("Inspeksi dan penanganan atap bocor"), BasePrice: intPtr(180000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("2-4 jam"), IsActive: true}},
		{"listrik", entity.Service{Name: "Instalasi Listrik", Slug: "instalasi-listrik", Description: strPtr("Pemasangan dan perbaikan instalasi listrik"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
		{"kunci", entity.Service{Name: "Perbaikan Kunci", Slug: "perbaikan-kunci", Description: strPtr("Perbaikan dan penggantian kunci rumah"), BasePrice: intPtr(100000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},
		{"kayu", entity.Service{Name: "Perbaikan Kusen", Slug: "perbaikan-kusen", Description: strPtr("Perbaikan kusen, pintu, dan furniture kayu"), BasePrice: intPtr(175000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("2-5 jam"), IsActive: true}},
		{"cat", entity.Service{Name: "Pengecatan Rumah", Slug: "pengecatan-rumah", Description: strPtr("Jasa pengecatan ruangan dan fasad rumah"), BasePrice: intPtr(250000), PriceUnit: strPtr("per ruangan"), EstimatedDuration: strPtr("4-8 jam"), IsActive: true}},
		{"kebun", entity.Service{Name: "Perawatan Taman", Slug: "perawatan-taman", Description: strPtr("Pembersihan dan perawatan taman rumah"), BasePrice: intPtr(125000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
	}
	for _, s := range seeds {
		var cat entity.Category
		if err := db.Where("slug = ?", s.CatSlug).First(&cat).Error; err != nil {
			return err
		}
		svc := s.Svc
		svc.CategoryID = cat.ID
		var existing entity.Service
		if err := db.Where("slug = ?", svc.Slug).Assign(svc).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── ARTICLES ─────────────────────────────────────────────────────────────────

func seedArticles(db *gorm.DB) error {
	now := time.Now()
	articles := []entity.Article{
		{
			Title: "Tips Merawat Rumah Sebelum Musim Hujan", Slug: "tips-merawat-rumah-sebelum-musim-hujan",
			Category: entity.ArticleCategoryTips, ThumbnailURL: strPtr("https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600"),
			Excerpt:         strPtr("Checklist sederhana untuk mencegah kerusakan atap, pipa, dan listrik."),
			ContentHTML:     "<p>Periksa talang, atap, sambungan pipa, dan panel listrik sebelum musim hujan tiba.</p>",
			ReadTimeMinutes: intPtr(3), Author: strPtr("Tim HandyDirect"),
			Tags: entity.JSONB(`["rumah","hujan","perawatan"]`), IsPublished: true, PublishedAt: &now,
		},
		{
			Title: "Panduan Memilih Tukang AC yang Terpercaya", Slug: "panduan-memilih-tukang-ac-terpercaya",
			Category: entity.ArticleCategoryGuide, ThumbnailURL: strPtr("https://images.unsplash.com/photo-1621905251918-48416bd8575a?w=600"),
			Excerpt:         strPtr("5 hal yang perlu diperhatikan saat memilih teknisi AC."),
			ContentHTML:     "<p>Pastikan teknisi AC memiliki sertifikasi, ulasan positif, dan jaminan garansi kerja.</p>",
			ReadTimeMinutes: intPtr(4), Author: strPtr("Tim HandyDirect"),
			Tags: entity.JSONB(`["AC","tips","teknisi"]`), IsPublished: true, PublishedAt: &now,
		},
		{
			Title: "Cara Aman Menangani Korsleting Listrik di Rumah", Slug: "cara-aman-tangani-korsleting-listrik",
			Category: entity.ArticleCategorySafety, ThumbnailURL: strPtr("https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600"),
			Excerpt:         strPtr("Langkah pertolongan pertama saat terjadi korsleting."),
			ContentHTML:     "<p>Matikan MCB utama, jangan sentuh perangkat basah, dan hubungi teknisi bersertifikat.</p>",
			ReadTimeMinutes: intPtr(5), Author: strPtr("Tim HandyDirect"),
			Tags: entity.JSONB(`["listrik","keamanan","korsleting"]`), IsPublished: true, PublishedAt: &now,
		},
	}
	for _, a := range articles {
		var existing entity.Article
		if err := db.Where("slug = ?", a.Slug).Assign(a).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── FAQS ─────────────────────────────────────────────────────────────────────

func seedFAQs(db *gorm.DB) error {
	faqs := []entity.FAQ{
		{Question: "Bagaimana cara memesan tukang?", Answer: "Pilih kategori jasa, pilih worker terdekat, jelaskan masalah, lalu buat pesanan. Tukang akan konfirmasi dalam 5-15 menit.", Category: entity.FAQCategoryGeneral, DisplayOrder: 1, IsActive: true},
		{Question: "Apakah ada jaminan kualitas pekerjaan?", Answer: "Ya. Semua tukang telah diverifikasi. Kami menyediakan garansi 7 hari untuk setiap pekerjaan yang selesai.", Category: entity.FAQCategoryGeneral, DisplayOrder: 2, IsActive: true},
		{Question: "Bagaimana cara membayar?", Answer: "Pembayaran dilakukan setelah pekerjaan selesai melalui cash, transfer bank, atau e-wallet.", Category: entity.FAQCategoryPayment, DisplayOrder: 1, IsActive: true},
		{Question: "Apakah bisa membatalkan pesanan?", Answer: "Bisa, selama tukang belum memulai perjalanan. Pembatalan setelah tukang berangkat dapat dikenakan biaya booking.", Category: entity.FAQCategoryCancellation, DisplayOrder: 1, IsActive: true},
		{Question: "Bagaimana cara melacak tukang?", Answer: "Setelah order diterima, Anda dapat melihat status real-time di halaman detail pesanan.", Category: entity.FAQCategoryTracking, DisplayOrder: 1, IsActive: true},
	}
	for _, f := range faqs {
		var existing entity.FAQ
		if err := db.Where("question = ?", f.Question).Assign(f).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── PROMOTIONS ───────────────────────────────────────────────────────────────

func seedPromotions(db *gorm.DB) error {
	now := time.Now()
	oneMonth := now.AddDate(0, 1, 0)
	twoMonth := now.AddDate(0, 2, 0)
	promos := []entity.Promotion{
		{
			Title: "Diskon 20% Jasa AC", Description: strPtr("Berlaku untuk servis dan pasang AC semua merk"),
			ImageURL: "https://images.unsplash.com/photo-1621905251918-48416bd8575a?w=800",
			CTALabel: strPtr("Klaim Sekarang"), DeepLink: strPtr("/promo/diskon-ac-20"),
			PromoCode: strPtr("AC20"), DisplayOrder: 1, IsActive: true,
			DiscountPercent: floatPtr(20), ValidFrom: &now, ValidUntil: &oneMonth,
		},
		{
			Title: "Gratis Booking Fee", Description: strPtr("Tidak ada booking fee untuk 3 order pertama Anda"),
			ImageURL: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800",
			CTALabel: strPtr("Gunakan Sekarang"), DeepLink: strPtr("/promo/gratis-booking"),
			PromoCode: strPtr("NEWUSER"), DisplayOrder: 2, IsActive: true,
			DiscountAmount: intPtr(2000), ValidFrom: &now, ValidUntil: &twoMonth,
		},
	}
	for _, p := range promos {
		var existing entity.Promotion
		if err := db.Where("promo_code = ?", p.PromoCode).Assign(p).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── NOTIFICATIONS ────────────────────────────────────────────────────────────

func seedNotifications(db *gorm.DB) error {
	type notifSeed struct {
		Email  string
		Notifs []entity.Notification
	}
	seeds := []notifSeed{
		{"siti.rahayu@gmail.com", []entity.Notification{
			{Type: entity.NotificationTypeOrder, Title: "Pesanan AC Selesai! ✅", Body: "Servis AC Anda oleh Budi Santoso telah selesai. Silakan beri ulasan.", DeepLink: strPtr("/orders/HD-20260501-001"), IsRead: false, Metadata: entity.JSONB(`{"order_number":"HD-20260501-001"}`)},
			{Type: entity.NotificationTypePromo, Title: "Promo Spesial Untukmu 🎉", Body: "Dapatkan diskon 20% untuk servis AC berikutnya dengan kode AC20.", DeepLink: strPtr("/promo/diskon-ac-20"), IsRead: true},
		}},
		{"dewi.kusuma@gmail.com", []entity.Notification{
			{Type: entity.NotificationTypeOrder, Title: "Pesanan Listrik Selesai! ✅", Body: "Instalasi listrik oleh Ahmad Fauzi telah selesai. Terima kasih!", DeepLink: strPtr("/orders/HD-20260510-002"), IsRead: false, Metadata: entity.JSONB(`{"order_number":"HD-20260510-002"}`)},
		}},
		{"andi.firmansyah@gmail.com", []entity.Notification{
			{Type: entity.NotificationTypeChat, Title: "Pesan Baru dari Hendra 💬", Body: "Hendra Wijaya mengirim pesan terkait perbaikan pipa Anda.", DeepLink: strPtr("/orders/HD-20260603-003/chat"), IsRead: false},
			{Type: entity.NotificationTypePurchase, Title: "Permintaan Pembelian Material 🛒", Body: "Hendra meminta persetujuan pembelian selang fleksibel Rp55.000.", DeepLink: strPtr("/orders/HD-20260603-003"), IsRead: false, Metadata: entity.JSONB(`{"order_number":"HD-20260603-003","amount":85000}`)},
		}},
		{"budi.santoso@situkang.id", []entity.Notification{
			{Type: entity.NotificationTypePayment, Title: "Pembayaran Diterima 💰", Body: "Pembayaran Rp272.000 dari Siti Rahayu untuk order HD-20260501-001 telah masuk.", DeepLink: strPtr("/wallet"), IsRead: true, Metadata: entity.JSONB(`{"amount":272000,"order_number":"HD-20260501-001"}`)},
		}},
		{"hendra.wijaya@situkang.id", []entity.Notification{
			{Type: entity.NotificationTypeOrder, Title: "Ada Order Baru! 🔔", Body: "Andi Firmansyah butuh bantuan perbaikan pipa bocor URGENT.", DeepLink: strPtr("/orders/HD-20260603-003"), IsRead: true},
		}},
	}

	for _, s := range seeds {
		var user entity.User
		if err := db.Where("email = ?", s.Email).First(&user).Error; err != nil {
			continue
		}
		for _, n := range s.Notifs {
			n.UserID = user.ID
			if err := db.Create(&n).Error; err != nil {
				continue
			}
		}
	}
	return nil
}

// ─── SHARED HELPERS ───────────────────────────────────────────────────────────

func strPtr(v string) *string   { return &v }
func intPtr(v int) *int         { return &v }
func floatPtr(v float64) *float64 { return &v }
func timePtr(v time.Time) *time.Time { return &v }
