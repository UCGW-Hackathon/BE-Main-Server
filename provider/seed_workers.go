package provider

import (
	"situkang/models/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedWorkers menyemai data tukang profesional lengkap dengan profil, layanan, dan wallet.
func SeedWorkers(db *gorm.DB) error {
	if err := seedWorkerUsers(db); err != nil {
		return err
	}
	if err := seedWorkerProfiles(db); err != nil {
		return err
	}
	if err := seedWorkerServicesData(db); err != nil {
		return err
	}
	if err := seedWorkerWallets(db); err != nil {
		return err
	}
	return nil
}

// ─── WORKER USERS ─────────────────────────────────────────────────────────────

// workerPasswordHash adalah bcrypt hash dari "Password123!" untuk semua seed worker
const workerPasswordHash = "$2a$12$WO4ic2KVlSAUH9dpJMyhX.qyZ3LTVI9Cgkxu6FRkC/EL4Pvy7HLxu"

// customerPasswordHash adalah bcrypt hash dari "Password123!" untuk semua seed user
const customerPasswordHash = "$2a$12$WO4ic2KVlSAUH9dpJMyhX.qyZ3LTVI9Cgkxu6FRkC/EL4Pvy7HLxu"

func seedWorkerUsers(db *gorm.DB) error {
	jakartaLat := -6.2088
	jakartaLng := 106.8456

	workers := []entity.User{
		{
			FullName:     "Budi Santoso",
			Email:        "budi.santoso@situkang.id",
			Phone:        "081234567801",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=budi.santoso"),
			Address:      strPtr("Jl. Kebon Jeruk No. 12, Jakarta Barat"),
			Latitude:     float64Ptr(jakartaLat + 0.01),
			Longitude:    float64Ptr(jakartaLng - 0.05),
			IsActive:     true,
		},
		{
			FullName:     "Ahmad Fauzi",
			Email:        "ahmad.fauzi@situkang.id",
			Phone:        "081234567802",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=ahmad.fauzi"),
			Address:      strPtr("Jl. Raya Pasar Minggu No. 45, Jakarta Selatan"),
			Latitude:     float64Ptr(jakartaLat - 0.03),
			Longitude:    float64Ptr(jakartaLng + 0.02),
			IsActive:     true,
		},
		{
			FullName:     "Hendra Wijaya",
			Email:        "hendra.wijaya@situkang.id",
			Phone:        "081234567803",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=hendra.wijaya"),
			Address:      strPtr("Jl. Mangga Besar No. 8, Jakarta Pusat"),
			Latitude:     float64Ptr(jakartaLat + 0.02),
			Longitude:    float64Ptr(jakartaLng + 0.01),
			IsActive:     true,
		},
		{
			FullName:     "Rizki Permana",
			Email:        "rizki.permana@situkang.id",
			Phone:        "081234567804",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=rizki.permana"),
			Address:      strPtr("Jl. Cipinang Besar No. 3, Jakarta Timur"),
			Latitude:     float64Ptr(jakartaLat + 0.05),
			Longitude:    float64Ptr(jakartaLng + 0.08),
			IsActive:     true,
		},
		{
			FullName:     "Doni Prasetyo",
			Email:        "doni.prasetyo@situkang.id",
			Phone:        "081234567805",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=doni.prasetyo"),
			Address:      strPtr("Jl. Kapuk No. 17, Jakarta Utara"),
			Latitude:     float64Ptr(jakartaLat + 0.12),
			Longitude:    float64Ptr(jakartaLng - 0.02),
			IsActive:     true,
		},
		{
			FullName:     "Eko Susanto",
			Email:        "eko.susanto@situkang.id",
			Phone:        "081234567806",
			PasswordHash: workerPasswordHash,
			Role:         entity.UserRoleWorker,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=eko.susanto"),
			Address:      strPtr("Jl. Cibubur Indah No. 22, Jakarta Timur"),
			Latitude:     float64Ptr(jakartaLat + 0.07),
			Longitude:    float64Ptr(jakartaLng + 0.10),
			IsActive:     true,
		},
	}

	for _, w := range workers {
		var existing entity.User
		if err := db.Where("email = ?", w.Email).Assign(w).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	// Seed customer users
	customers := []entity.User{
		{
			FullName:     "Siti Rahayu",
			Email:        "siti.rahayu@gmail.com",
			Phone:        "081298765401",
			PasswordHash: customerPasswordHash,
			Role:         entity.UserRoleUser,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=siti.rahayu"),
			Address:      strPtr("Jl. Sudirman No. 55, Jakarta Pusat"),
			Latitude:     float64Ptr(-6.2100),
			Longitude:    float64Ptr(106.8230),
			IsActive:     true,
		},
		{
			FullName:     "Dewi Kusuma",
			Email:        "dewi.kusuma@gmail.com",
			Phone:        "081298765402",
			PasswordHash: customerPasswordHash,
			Role:         entity.UserRoleUser,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=dewi.kusuma"),
			Address:      strPtr("Jl. Kemang Raya No. 10, Jakarta Selatan"),
			Latitude:     float64Ptr(-6.2600),
			Longitude:    float64Ptr(106.8150),
			IsActive:     true,
		},
		{
			FullName:     "Andi Firmansyah",
			Email:        "andi.firmansyah@gmail.com",
			Phone:        "081298765403",
			PasswordHash: customerPasswordHash,
			Role:         entity.UserRoleUser,
			AvatarURL:    strPtr("https://i.pravatar.cc/150?u=andi.firmansyah"),
			Address:      strPtr("Jl. Pondok Indah No. 33, Jakarta Selatan"),
			Latitude:     float64Ptr(-6.2700),
			Longitude:    float64Ptr(106.7900),
			IsActive:     true,
		},
	}

	for _, c := range customers {
		var existing entity.User
		if err := db.Where("email = ?", c.Email).Assign(c).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	return nil
}

// ─── WORKER PROFILES ──────────────────────────────────────────────────────────

func seedWorkerProfiles(db *gorm.DB) error {
	verifiedAt := time.Now().AddDate(0, -6, 0)

	profiles := []struct {
		Email   string
		Profile entity.WorkerProfile
	}{
		{
			Email: "budi.santoso@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Teknisi AC & Pendingin"),
				Bio:                strPtr("Teknisi AC berpengalaman 10 tahun. Melayani servis, pemasangan, dan isi freon semua merk AC. Bergaransi 30 hari."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1621905251918-48416bd8575a?w=800"),
				VerificationStatus: entity.VerificationStatusVerified,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/budi-santoso.jpg"),
				CertificateURLs:    entity.JSONB(`["https://cdn.handydirect.id/certs/ac-technician-budi.pdf"]`),
				BasePrice:          intPtr(150000),
				PriceUnit:          strPtr("per kunjungan"),
				BookingFee:         2000,
				RatingAvg:          4.8,
				TotalReviews:       47,
				CompletedJobs:      53,
				IsAvailable:        true,
				VerifiedAt:         &verifiedAt,
			},
		},
		{
			Email: "ahmad.fauzi@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Teknisi Listrik & Instalasi"),
				Bio:                strPtr("Tukang listrik bersertifikat PLN dengan pengalaman 8 tahun. Spesialis instalasi panel, stop kontak, dan perbaikan korsleting."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800"),
				VerificationStatus: entity.VerificationStatusVerified,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/ahmad-fauzi.jpg"),
				CertificateURLs:    entity.JSONB(`["https://cdn.handydirect.id/certs/listrik-pln-ahmad.pdf","https://cdn.handydirect.id/certs/k3-listrik-ahmad.pdf"]`),
				BasePrice:          intPtr(175000),
				PriceUnit:          strPtr("per kunjungan"),
				BookingFee:         2000,
				RatingAvg:          4.9,
				TotalReviews:       62,
				CompletedJobs:      70,
				IsAvailable:        true,
				VerifiedAt:         &verifiedAt,
			},
		},
		{
			Email: "hendra.wijaya@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Tukang Pipa & Sanitasi"),
				Bio:                strPtr("Ahli pipa dan sanitasi rumah tangga 7 tahun. Berpengalaman menangani kebocoran, saluran mampet, dan instalasi pipa baru."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1504328345606-18bbc8c9d7d1?w=800"),
				VerificationStatus: entity.VerificationStatusVerified,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/hendra-wijaya.jpg"),
				CertificateURLs:    entity.JSONB(`["https://cdn.handydirect.id/certs/sanitasi-hendra.pdf"]`),
				BasePrice:          intPtr(120000),
				PriceUnit:          strPtr("per kunjungan"),
				BookingFee:         2000,
				RatingAvg:          4.7,
				TotalReviews:       38,
				CompletedJobs:      44,
				IsAvailable:        true,
				VerifiedAt:         &verifiedAt,
			},
		},
		{
			Email: "rizki.permana@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Tukang Kayu & Furniture"),
				Bio:                strPtr("Ahli kayu dan furniture 12 tahun. Spesialis pintu, kusen, lemari, dan renovasi interior berbahan kayu."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1588854337115-1c67d9247e4d?w=800"),
				VerificationStatus: entity.VerificationStatusVerified,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/rizki-permana.jpg"),
				CertificateURLs:    entity.JSONB(`["https://cdn.handydirect.id/certs/furniture-rizki.pdf"]`),
				BasePrice:          intPtr(200000),
				PriceUnit:          strPtr("per kunjungan"),
				BookingFee:         2000,
				RatingAvg:          4.6,
				TotalReviews:       29,
				CompletedJobs:      35,
				IsAvailable:        true,
				VerifiedAt:         &verifiedAt,
			},
		},
		{
			Email: "doni.prasetyo@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Tukang Cat & Dekorasi"),
				Bio:                strPtr("Pelukis dan pengcat profesional 9 tahun. Interior, eksterior, waterproofing, dan dekorasi dinding rumah dan apartemen."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1562259929-b4e1fd3aef09?w=800"),
				VerificationStatus: entity.VerificationStatusVerified,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/doni-prasetyo.jpg"),
				CertificateURLs:    entity.JSONB(`[]`),
				BasePrice:          intPtr(250000),
				PriceUnit:          strPtr("per ruangan"),
				BookingFee:         2000,
				RatingAvg:          4.5,
				TotalReviews:       21,
				CompletedJobs:      26,
				IsAvailable:        true,
				VerifiedAt:         &verifiedAt,
			},
		},
		{
			Email: "eko.susanto@situkang.id",
			Profile: entity.WorkerProfile{
				Specialization:     strPtr("Tukang Atap & Bocor"),
				Bio:                strPtr("Spesialis atap bocor dan waterproofing 6 tahun. Berpengalaman menangani genteng, spandek, bitumen, dan beton."),
				CoverPhotoURL:      strPtr("https://images.unsplash.com/photo-1590859808308-3d2d9c515b1a?w=800"),
				VerificationStatus: entity.VerificationStatusPending,
				IDCardURL:          strPtr("https://cdn.handydirect.id/id-cards/eko-susanto.jpg"),
				CertificateURLs:    entity.JSONB(`[]`),
				BasePrice:          intPtr(180000),
				PriceUnit:          strPtr("per kunjungan"),
				BookingFee:         2000,
				RatingAvg:          4.3,
				TotalReviews:       14,
				CompletedJobs:      17,
				IsAvailable:        true,
			},
		},
	}

	for _, p := range profiles {
		var user entity.User
		if err := db.Where("email = ?", p.Email).First(&user).Error; err != nil {
			continue
		}
		p.Profile.UserID = user.ID
		var existing entity.WorkerProfile
		if err := db.Where("user_id = ?", user.ID).Assign(p.Profile).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── WORKER SERVICES ──────────────────────────────────────────────────────────

func seedWorkerServicesData(db *gorm.DB) error {
	type workerServiceSeed struct {
		WorkerEmail string
		ServiceSlug string
		CustomPrice *int
	}

	seeds := []workerServiceSeed{
		// Budi: AC specialist
		{"budi.santoso@situkang.id", "servis-ac", intPtr(150000)},
		{"budi.santoso@situkang.id", "pemasangan-ac", intPtr(350000)},

		// Ahmad: Electrician
		{"ahmad.fauzi@situkang.id", "instalasi-listrik", intPtr(175000)},
		{"ahmad.fauzi@situkang.id", "perbaikan-listrik", intPtr(150000)},

		// Hendra: Plumber
		{"hendra.wijaya@situkang.id", "perbaikan-pipa", intPtr(120000)},
		{"hendra.wijaya@situkang.id", "instalasi-pipa-baru", intPtr(200000)},

		// Rizki: Woodworker
		{"rizki.permana@situkang.id", "perbaikan-kusen", intPtr(200000)},
		{"rizki.permana@situkang.id", "instalasi-pintu", intPtr(350000)},

		// Doni: Painter
		{"doni.prasetyo@situkang.id", "pengecatan-rumah", intPtr(300000)},
		{"doni.prasetyo@situkang.id", "pengecatan-eksterior", intPtr(400000)},

		// Eko: Roofer
		{"eko.susanto@situkang.id", "perbaikan-atap-bocor", intPtr(180000)},
	}

	for _, s := range seeds {
		var user entity.User
		if err := db.Where("email = ?", s.WorkerEmail).First(&user).Error; err != nil {
			continue
		}
		var service entity.Service
		if err := db.Where("slug = ?", s.ServiceSlug).First(&service).Error; err != nil {
			continue
		}
		ws := entity.WorkerService{
			WorkerID:    user.ID,
			ServiceID:   service.ID,
			CustomPrice: s.CustomPrice,
			IsActive:    true,
		}
		var existing entity.WorkerService
		if err := db.Where("worker_id = ? AND service_id = ?", user.ID, service.ID).
			Assign(ws).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── WORKER WALLETS ───────────────────────────────────────────────────────────

func seedWorkerWallets(db *gorm.DB) error {
	type walletSeed struct {
		WorkerEmail    string
		Balance        int64
		TotalEarnings  int64
		TotalWithdrawn int64
	}

	seeds := []walletSeed{
		{"budi.santoso@situkang.id", 1850000, 8750000, 6900000},
		{"ahmad.fauzi@situkang.id", 2350000, 12500000, 10150000},
		{"hendra.wijaya@situkang.id", 980000, 5500000, 4520000},
		{"rizki.permana@situkang.id", 1450000, 7000000, 5550000},
		{"doni.prasetyo@situkang.id", 760000, 3900000, 3140000},
		{"eko.susanto@situkang.id", 340000, 2100000, 1760000},
	}

	for _, s := range seeds {
		var user entity.User
		if err := db.Where("email = ?", s.WorkerEmail).First(&user).Error; err != nil {
			continue
		}
		wallet := entity.WorkerWallet{
			WorkerID:       user.ID,
			Balance:        s.Balance,
			TotalEarnings:  s.TotalEarnings,
			TotalWithdrawn: s.TotalWithdrawn,
			PendingEarnings: 0,
			IsActive:       true,
		}
		var existing entity.WorkerWallet
		if err := db.Where("worker_id = ?", user.ID).Assign(wallet).FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── HELPERS ──────────────────────────────────────────────────────────────────

func float64Ptr(v float64) *float64 { return &v }

// GetWorkerUUID mengambil UUID worker berdasarkan email (helper untuk seed lain).
func GetWorkerUUID(db *gorm.DB, email string) (uuid.UUID, error) {
	var user entity.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}
