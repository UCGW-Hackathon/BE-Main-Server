package provider

import (
	"time"

	"situkang/models/entity"

	"gorm.io/gorm"
)

func seedReferenceData(db *gorm.DB) error {
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

	for _, category := range categories {
		var existing entity.Category
		if err := db.Where("slug = ?", category.Slug).Assign(category).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	serviceSeeds := []struct {
		CategorySlug string
		Service      entity.Service
	}{
		{"ac", entity.Service{Name: "Servis AC", Slug: "servis-ac", Description: strPtr("Cuci, tambah freon, dan perbaikan AC"), IconURL: strPtr("https://cdn.handydirect.id/icons/servis_ac.png"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
		{"pipa", entity.Service{Name: "Perbaikan Pipa", Slug: "perbaikan-pipa", Description: strPtr("Perbaikan pipa bocor dan saluran air"), IconURL: strPtr("https://cdn.handydirect.id/icons/perbaikan_pipa.png"), BasePrice: intPtr(120000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},
		{"atap", entity.Service{Name: "Perbaikan Atap Bocor", Slug: "perbaikan-atap-bocor", Description: strPtr("Inspeksi dan penanganan atap bocor"), IconURL: strPtr("https://cdn.handydirect.id/icons/atap.png"), BasePrice: intPtr(180000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("2-4 jam"), IsActive: true}},
		{"listrik", entity.Service{Name: "Instalasi Listrik", Slug: "instalasi-listrik", Description: strPtr("Pemasangan dan perbaikan instalasi listrik"), IconURL: strPtr("https://cdn.handydirect.id/icons/instalasi_listrik.png"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
		{"kunci", entity.Service{Name: "Perbaikan Kunci", Slug: "perbaikan-kunci", Description: strPtr("Perbaikan dan penggantian kunci rumah"), IconURL: strPtr("https://cdn.handydirect.id/icons/kunci.png"), BasePrice: intPtr(100000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},
		{"kayu", entity.Service{Name: "Perbaikan Kusen", Slug: "perbaikan-kusen", Description: strPtr("Perbaikan kusen, pintu, dan furniture kayu"), IconURL: strPtr("https://cdn.handydirect.id/icons/kayu.png"), BasePrice: intPtr(175000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("2-5 jam"), IsActive: true}},
		{"cat", entity.Service{Name: "Pengecatan Rumah", Slug: "pengecatan-rumah", Description: strPtr("Jasa pengecatan ruangan dan fasad rumah"), IconURL: strPtr("https://cdn.handydirect.id/icons/cat.png"), BasePrice: intPtr(250000), PriceUnit: strPtr("per ruangan"), EstimatedDuration: strPtr("4-8 jam"), IsActive: true}},
		{"kebun", entity.Service{Name: "Perawatan Taman", Slug: "perawatan-taman", Description: strPtr("Pembersihan dan perawatan taman rumah"), IconURL: strPtr("https://cdn.handydirect.id/icons/kebun.png"), BasePrice: intPtr(125000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},
	}

	for _, seed := range serviceSeeds {
		var category entity.Category
		if err := db.Where("slug = ?", seed.CategorySlug).First(&category).Error; err != nil {
			return err
		}
		service := seed.Service
		service.CategoryID = category.ID
		var existing entity.Service
		if err := db.Where("slug = ?", service.Slug).Assign(service).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	now := time.Now()
	article := entity.Article{
		Title:           "Tips Merawat Rumah Sebelum Musim Hujan",
		Slug:            "tips-merawat-rumah-sebelum-musim-hujan",
		Category:        entity.ArticleCategoryTips,
		ThumbnailURL:    strPtr("https://cdn.handydirect.id/articles/tips_rumah_hujan.jpg"),
		Excerpt:         strPtr("Checklist sederhana untuk mencegah kerusakan atap, pipa, dan listrik."),
		ContentHTML:     "<p>Periksa talang, atap, sambungan pipa, dan panel listrik sebelum musim hujan tiba.</p>",
		ReadTimeMinutes: intPtr(3),
		Author:          strPtr("Tim HandyDirect"),
		Tags:            entity.JSONB(`["rumah","hujan","perawatan"]`),
		IsPublished:     true,
		PublishedAt:     &now,
	}
	var existingArticle entity.Article
	if err := db.Where("slug = ?", article.Slug).Assign(article).FirstOrCreate(&existingArticle).Error; err != nil {
		return err
	}

	faq := entity.FAQ{
		Question:     "Bagaimana cara memesan tukang?",
		Answer:       "Pilih kategori jasa, pilih worker, jelaskan masalah, lalu buat pesanan.",
		Category:     entity.FAQCategoryGeneral,
		DisplayOrder: 1,
		IsActive:     true,
	}
	var existingFAQ entity.FAQ
	if err := db.Where("question = ?", faq.Question).Assign(faq).FirstOrCreate(&existingFAQ).Error; err != nil {
		return err
	}

	promo := entity.Promotion{
		Title:           "Diskon 20% Jasa AC",
		Description:     strPtr("Berlaku untuk layanan AC pilihan"),
		ImageURL:        "https://cdn.handydirect.id/promos/diskon_ac.jpg",
		CTALabel:        strPtr("Klaim Sekarang"),
		DeepLink:        strPtr("/promo/diskon-ac-20"),
		PromoCode:       strPtr("AC20"),
		DisplayOrder:    1,
		IsActive:        true,
		DiscountPercent: floatPtr(20),
		ValidFrom:       &now,
		ValidUntil:      timePtr(now.AddDate(0, 1, 0)),
	}
	var existingPromo entity.Promotion
	return db.Where("promo_code = ?", promo.PromoCode).Assign(promo).FirstOrCreate(&existingPromo).Error
}

func strPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}
