package provider

import (
	"situkang/models/entity"

	"gorm.io/gorm"
)

// SeedExtendedServices menambahkan layanan tambahan di setiap kategori.
func SeedExtendedServices(db *gorm.DB) error {
	type svcSeed struct {
		CategorySlug string
		Service      entity.Service
	}

	seeds := []svcSeed{
		// AC
		{"ac", entity.Service{Name: "Pemasangan AC", Slug: "pemasangan-ac", Description: strPtr("Pemasangan unit AC baru termasuk instalasi pipa freon"), IconURL: strPtr("https://cdn.handydirect.id/icons/pasang_ac.png"), BasePrice: intPtr(350000), PriceUnit: strPtr("per unit"), EstimatedDuration: strPtr("2-4 jam"), IsActive: true}},
		{"ac", entity.Service{Name: "Isi Freon AC", Slug: "isi-freon-ac", Description: strPtr("Penambahan freon R22/R32 untuk AC kurang dingin"), IconURL: strPtr("https://cdn.handydirect.id/icons/freon_ac.png"), BasePrice: intPtr(120000), PriceUnit: strPtr("per unit"), EstimatedDuration: strPtr("30-60 menit"), IsActive: true}},

		// Listrik
		{"listrik", entity.Service{Name: "Perbaikan Listrik", Slug: "perbaikan-listrik", Description: strPtr("Perbaikan korsleting, MCB trip, dan masalah listrik lainnya"), IconURL: strPtr("https://cdn.handydirect.id/icons/perbaikan_listrik.png"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},
		{"listrik", entity.Service{Name: "Pasang Stop Kontak", Slug: "pasang-stop-kontak", Description: strPtr("Pemasangan stop kontak dan saklar tambahan"), IconURL: strPtr("https://cdn.handydirect.id/icons/stop_kontak.png"), BasePrice: intPtr(80000), PriceUnit: strPtr("per titik"), EstimatedDuration: strPtr("30-60 menit"), IsActive: true}},

		// Pipa
		{"pipa", entity.Service{Name: "Instalasi Pipa Baru", Slug: "instalasi-pipa-baru", Description: strPtr("Instalasi jalur pipa air bersih dan kotor untuk renovasi"), IconURL: strPtr("https://cdn.handydirect.id/icons/pipa_baru.png"), BasePrice: intPtr(200000), PriceUnit: strPtr("per meter"), EstimatedDuration: strPtr("2-5 jam"), IsActive: true}},
		{"pipa", entity.Service{Name: "Saluran Mampet", Slug: "saluran-mampet", Description: strPtr("Pembersihan saluran air dan toilet yang tersumbat"), IconURL: strPtr("https://cdn.handydirect.id/icons/saluran_mampet.png"), BasePrice: intPtr(100000), PriceUnit: strPtr("per titik"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},

		// Atap
		{"atap", entity.Service{Name: "Waterproofing Atap", Slug: "waterproofing-atap", Description: strPtr("Aplikasi waterproofing untuk mencegah kebocoran atap beton"), IconURL: strPtr("https://cdn.handydirect.id/icons/waterproofing.png"), BasePrice: intPtr(300000), PriceUnit: strPtr("per m2"), EstimatedDuration: strPtr("4-8 jam"), IsActive: true}},

		// Kayu
		{"kayu", entity.Service{Name: "Instalasi Pintu", Slug: "instalasi-pintu", Description: strPtr("Pemasangan dan setting pintu baru beserta engsel dan kunci"), IconURL: strPtr("https://cdn.handydirect.id/icons/pintu.png"), BasePrice: intPtr(350000), PriceUnit: strPtr("per daun pintu"), EstimatedDuration: strPtr("2-4 jam"), IsActive: true}},
		{"kayu", entity.Service{Name: "Perbaikan Furniture", Slug: "perbaikan-furniture", Description: strPtr("Perbaikan lemari, meja, dan furniture kayu lainnya"), IconURL: strPtr("https://cdn.handydirect.id/icons/furniture.png"), BasePrice: intPtr(150000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-3 jam"), IsActive: true}},

		// Cat
		{"cat", entity.Service{Name: "Pengecatan Eksterior", Slug: "pengecatan-eksterior", Description: strPtr("Pengecatan fasad dan eksterior bangunan"), IconURL: strPtr("https://cdn.handydirect.id/icons/cat_eksterior.png"), BasePrice: intPtr(400000), PriceUnit: strPtr("per m2"), EstimatedDuration: strPtr("4-8 jam"), IsActive: true}},
		{"cat", entity.Service{Name: "Pengecatan Pagar", Slug: "pengecatan-pagar", Description: strPtr("Cat ulang pagar besi dan tembok pagar"), IconURL: strPtr("https://cdn.handydirect.id/icons/cat_pagar.png"), BasePrice: intPtr(150000), PriceUnit: strPtr("per meter"), EstimatedDuration: strPtr("2-4 jam"), IsActive: true}},

		// Kebun
		{"kebun", entity.Service{Name: "Pangkas Rumput", Slug: "pangkas-rumput", Description: strPtr("Pemotongan dan perapian rumput taman"), IconURL: strPtr("https://cdn.handydirect.id/icons/pangkas_rumput.png"), BasePrice: intPtr(100000), PriceUnit: strPtr("per kunjungan"), EstimatedDuration: strPtr("1-2 jam"), IsActive: true}},

		// Kunci
		{"kunci", entity.Service{Name: "Ganti Kunci Pintu", Slug: "ganti-kunci-pintu", Description: strPtr("Penggantian handle dan silinder kunci pintu"), IconURL: strPtr("https://cdn.handydirect.id/icons/ganti_kunci.png"), BasePrice: intPtr(120000), PriceUnit: strPtr("per kunci"), EstimatedDuration: strPtr("30-60 menit"), IsActive: true}},
	}

	for _, s := range seeds {
		var category entity.Category
		if err := db.Where("slug = ?", s.CategorySlug).First(&category).Error; err != nil {
			continue
		}
		svc := s.Service
		svc.CategoryID = category.ID
		var existing entity.Service
		if err := db.Where("slug = ?", svc.Slug).Assign(svc).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}
