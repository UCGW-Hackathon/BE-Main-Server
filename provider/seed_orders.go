package provider

import (
	"situkang/models/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedOrdersAndHistory menyemai data order lengkap: timeline, chat, purchases, reviews, invoices, payments.
func SeedOrdersAndHistory(db *gorm.DB) error {
	// Fetch user IDs
	type UserRef struct {
		Email string
		ID    uuid.UUID
	}
	emailRefs := []string{
		"siti.rahayu@gmail.com",
		"dewi.kusuma@gmail.com",
		"andi.firmansyah@gmail.com",
		"budi.santoso@situkang.id",
		"ahmad.fauzi@situkang.id",
		"hendra.wijaya@situkang.id",
		"rizki.permana@situkang.id",
		"doni.prasetyo@situkang.id",
	}
	userMap := make(map[string]uuid.UUID)
	for _, email := range emailRefs {
		var u entity.User
		if err := db.Where("email = ?", email).First(&u).Error; err != nil {
			continue
		}
		userMap[email] = u.ID
	}

	// Fetch service & category IDs
	serviceMap := make(map[string]uuid.UUID)
	categoryMap := make(map[string]uuid.UUID)
	var services []entity.Service
	db.Find(&services)
	for _, s := range services {
		serviceMap[s.Slug] = s.ID
		categoryMap[s.Slug] = s.CategoryID
	}

	now := time.Now()
	ago := func(d time.Duration) time.Time { return now.Add(-d) }
	ptr := func(t time.Time) *time.Time { return &t }
	iptr := func(i int) *int { return &i }

	type OrderSeed struct {
		Order     entity.Order
		Photos    []entity.OrderPhoto
		Timeline  []entity.OrderTimeline
		Chats     []entity.ChatMessage
		Purchases []entity.Purchase
		Review    *entity.Review
		ReviewTags []string
	}

	seeds := []OrderSeed{
		// ── ORDER 1: Completed AC Service ─────────────────────────────────────
		{
			Order: entity.Order{
				OrderNumber:     "HD-20260501-001",
				UserID:          userMap["siti.rahayu@gmail.com"],
				WorkerID:        userMap["budi.santoso@situkang.id"],
				ServiceID:       serviceMap["servis-ac"],
				CategoryID:      categoryMap["servis-ac"],
				Title:           "AC kamar tidak dingin, perlu servis",
				Description:     "AC 1 PK Sharp sudah tidak dingin sejak 2 minggu. Sudah dicoba restart tapi tetap tidak dingin. Mungkin perlu tambah freon.",
				Status:          entity.OrderStatusCompleted,
				Urgency:         entity.OrderUrgencyNormal,
				LocationAddress: "Jl. Sudirman No. 55, Jakarta Pusat",
				LocationDetail:  strPtr("Lantai 3, kamar tidur utama"),
				LocationLat:     -6.2100,
				LocationLng:     106.8230,
				PreferredDate:   ptr(ago(10 * 24 * time.Hour)),
				Notes:           strPtr("Tolong bawa alat lengkap"),
				BookingFee:      2000,
				BaseServiceFee:  iptr(150000),
				TotalMaterialCost: 120000,
				GrandTotal:      iptr(272000),
				AcceptedAt:      ptr(ago(10*24*time.Hour + 1*time.Hour)),
				StartedAt:       ptr(ago(10*24*time.Hour - 1*time.Hour)),
				CompletedAt:     ptr(ago(10*24*time.Hour - 3*time.Hour)),
			},
			Photos: []entity.OrderPhoto{
				{PhotoURL: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600", Caption: strPtr("AC Sharp 1 PK tampak depan"), DisplayOrder: 0},
				{PhotoURL: "https://images.unsplash.com/photo-1584622650111-993a426fbf0a?w=600", Caption: strPtr("Remote AC"), DisplayOrder: 1},
			},
			Timeline: []entity.OrderTimeline{
				{Event: "order_created", Label: "Pesanan Dibuat", Description: strPtr("Pesanan berhasil dibuat oleh pelanggan"), ActorType: strPtr("user")},
				{Event: "order_accepted", Label: "Pesanan Diterima", Description: strPtr("Tukang menerima pesanan dan akan segera berangkat"), ActorType: strPtr("worker")},
				{Event: "worker_on_the_way", Label: "Tukang Berangkat", Description: strPtr("Tukang sedang dalam perjalanan ke lokasi"), ActorType: strPtr("worker")},
				{Event: "worker_arrived", Label: "Tukang Tiba", Description: strPtr("Tukang telah tiba di lokasi"), ActorType: strPtr("worker")},
				{Event: "work_started", Label: "Pekerjaan Dimulai", Description: strPtr("Teknisi mulai memeriksa unit AC"), ActorType: strPtr("worker")},
				{Event: "work_completed", Label: "Pekerjaan Selesai", Description: strPtr("Servis AC selesai, freon ditambah, AC sudah dingin kembali"), ActorType: strPtr("worker")},
			},
			Chats: []entity.ChatMessage{
				{SenderID: userMap["siti.rahayu@gmail.com"], SenderType: "user", Content: strPtr("Pak Budi, AC saya di kamar sudah tidak dingin sejak 2 minggu. Kira-kira kenapa ya?"), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["budi.santoso@situkang.id"], SenderType: "worker", Content: strPtr("Halo Bu Siti, kemungkinan freon kurang atau filter kotor. Nanti saya cek dulu ya Bu waktu datang."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["siti.rahayu@gmail.com"], SenderType: "user", Content: strPtr("Baik Pak, ditunggu ya. AC merk Sharp 1 PK, sudah 3 tahun."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["budi.santoso@situkang.id"], SenderType: "worker", Content: strPtr("Siap Bu, saya sudah tiba di depan ya."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["budi.santoso@situkang.id"], SenderType: "worker", Content: strPtr("Bu Siti, sudah saya cek. Freon R22 berkurang dan filter sangat kotor. Perlu tambah freon 1/2 kg dan cuci filter. Total material Rp120.000. Boleh dilanjut?"), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["siti.rahayu@gmail.com"], SenderType: "user", Content: strPtr("Boleh Pak, lanjut saja."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["budi.santoso@situkang.id"], SenderType: "worker", Content: strPtr("Alhamdulillah selesai Bu. AC sudah dingin kembali. Terima kasih sudah mempercayai saya."), MessageType: entity.MessageTypeText, IsRead: true},
			},
			Purchases: []entity.Purchase{
				{
					ItemName: "Freon R22 1/2 kg", Category: entity.PurchaseCategoryMaterial,
					Quantity: 0.5, Unit: "kg", UnitPrice: 120000, TotalPrice: 60000,
					Reason: strPtr("Freon AC kurang, perlu diisi agar AC kembali dingin"),
					ReceiptPhotoURL: strPtr("https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=400"),
					Status: entity.PurchaseStatusApproved, Confidence: float64Ptr(0.95),
					AIExplanation: strPtr("Item freon R22 sesuai dengan kebutuhan servis AC yang diorder pelanggan."),
				},
				{
					ItemName: "Cairan Cuci AC", Category: entity.PurchaseCategoryMaterial,
					Quantity: 1, Unit: "botol", UnitPrice: 60000, TotalPrice: 60000,
					Reason: strPtr("Filter AC sangat kotor, perlu cairan khusus untuk pembersihan"),
					Status: entity.PurchaseStatusApproved, Confidence: float64Ptr(0.92),
					AIExplanation: strPtr("Pembelian cairan cuci AC wajar dan relevan dengan pekerjaan servis AC."),
				},
			},
			Review: &entity.Review{
				Rating:  5,
				Comment: strPtr("Pak Budi sangat profesional! Datang tepat waktu, AC langsung dingin setelah dikerjakan. Penjelasannya juga detail. Sangat puas!"),
			},
			ReviewTags: []string{"tepat waktu", "profesional", "rapi", "ramah"},
		},

		// ── ORDER 2: Completed Electrical ─────────────────────────────────────
		{
			Order: entity.Order{
				OrderNumber:     "HD-20260510-002",
				UserID:          userMap["dewi.kusuma@gmail.com"],
				WorkerID:        userMap["ahmad.fauzi@situkang.id"],
				ServiceID:       serviceMap["instalasi-listrik"],
				CategoryID:      categoryMap["instalasi-listrik"],
				Title:           "Instalasi stop kontak tambahan di ruang kerja",
				Description:     "Butuh tambah 4 titik stop kontak di ruang kerja rumah. Ada 1 saklar rusak juga perlu diganti.",
				Status:          entity.OrderStatusCompleted,
				Urgency:         entity.OrderUrgencyNormal,
				LocationAddress: "Jl. Kemang Raya No. 10, Jakarta Selatan",
				LocationDetail:  strPtr("Ruang kerja lantai 1"),
				LocationLat:     -6.2600,
				LocationLng:     106.8150,
				Notes:           strPtr("Kabel sudah ada, tinggal pasang"),
				BookingFee:      2000,
				BaseServiceFee:  iptr(175000),
				TotalMaterialCost: 285000,
				GrandTotal:      iptr(462000),
				AcceptedAt:      ptr(ago(5 * 24 * time.Hour)),
				StartedAt:       ptr(ago(5*24*time.Hour - 2*time.Hour)),
				CompletedAt:     ptr(ago(5*24*time.Hour - 5*time.Hour)),
			},
			Photos: []entity.OrderPhoto{
				{PhotoURL: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600", Caption: strPtr("Panel listrik lama"), DisplayOrder: 0},
			},
			Timeline: []entity.OrderTimeline{
				{Event: "order_created", Label: "Pesanan Dibuat", ActorType: strPtr("user")},
				{Event: "order_accepted", Label: "Pesanan Diterima", ActorType: strPtr("worker")},
				{Event: "worker_arrived", Label: "Tukang Tiba", ActorType: strPtr("worker")},
				{Event: "work_started", Label: "Pekerjaan Dimulai", ActorType: strPtr("worker")},
				{Event: "work_completed", Label: "Pekerjaan Selesai", Description: strPtr("4 stop kontak dan 1 saklar berhasil dipasang"), ActorType: strPtr("worker")},
			},
			Chats: []entity.ChatMessage{
				{SenderID: userMap["dewi.kusuma@gmail.com"], SenderType: "user", Content: strPtr("Pak Ahmad, saya butuh tambah stop kontak 4 titik di ruang kerja. Bisa hari ini?"), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["ahmad.fauzi@situkang.id"], SenderType: "worker", Content: strPtr("Bisa Bu Dewi, saya sudah pesan. Kira-kira 2 jam pengerjaan. Kabel sudah ada atau perlu beli?"), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["dewi.kusuma@gmail.com"], SenderType: "user", Content: strPtr("Kabel ada, tapi stop kontaknya belum ada. Tolong sekalian belikan ya Pak."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["ahmad.fauzi@situkang.id"], SenderType: "worker", Content: strPtr("Siap Bu. Saya beli stop kontak Broco 4 titik + 1 saklar = Rp285.000. Saya kirimkan nota fotonya ya."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["dewi.kusuma@gmail.com"], SenderType: "user", Content: strPtr("Ok Pak, diapprove."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["ahmad.fauzi@situkang.id"], SenderType: "worker", Content: strPtr("Selesai Bu Dewi. Semua stop kontak dan saklar sudah terpasang dan berfungsi dengan baik."), MessageType: entity.MessageTypeText, IsRead: true},
			},
			Purchases: []entity.Purchase{
				{
					ItemName: "Stop Kontak Broco 3 lubang", Category: entity.PurchaseCategoryMaterial,
					Quantity: 4, Unit: "pcs", UnitPrice: 55000, TotalPrice: 220000,
					Reason:          strPtr("Stop kontak 4 titik sesuai permintaan pelanggan"),
					ReceiptPhotoURL: strPtr("https://images.unsplash.com/photo-1565193566173-7a0ee3dbe261?w=400"),
					Status: entity.PurchaseStatusApproved, Confidence: float64Ptr(0.97),
				},
				{
					ItemName: "Saklar Tunggal Broco", Category: entity.PurchaseCategoryMaterial,
					Quantity: 1, Unit: "pcs", UnitPrice: 65000, TotalPrice: 65000,
					Reason: strPtr("1 saklar rusak perlu diganti"),
					Status: entity.PurchaseStatusApproved, Confidence: float64Ptr(0.96),
				},
			},
			Review: &entity.Review{
				Rating:  5,
				Comment: strPtr("Pak Ahmad sangat ahli di bidang listrik. Kerjanya bersih, tidak ada kabel berantakan. Highly recommended!"),
			},
			ReviewTags: []string{"profesional", "rapi", "cepat", "bersih"},
		},

		// ── ORDER 3: In Progress Plumbing ──────────────────────────────────────
		{
			Order: entity.Order{
				OrderNumber:     "HD-20260603-003",
				UserID:          userMap["andi.firmansyah@gmail.com"],
				WorkerID:        userMap["hendra.wijaya@situkang.id"],
				ServiceID:       serviceMap["perbaikan-pipa"],
				CategoryID:      categoryMap["perbaikan-pipa"],
				Title:           "Pipa bocor di dapur, air mengalir terus",
				Description:     "Ada kebocoran di bawah wastafel dapur. Air terus menetes dan sudah basah lemari bawah. Perlu segera diperbaiki.",
				Status:          entity.OrderStatusInProgress,
				Urgency:         entity.OrderUrgencyUrgent,
				LocationAddress: "Jl. Pondok Indah No. 33, Jakarta Selatan",
				LocationDetail:  strPtr("Dapur belakang, wastafel"),
				LocationLat:     -6.2700,
				LocationLng:     106.7900,
				BookingFee:      2000,
				BaseServiceFee:  iptr(120000),
				TotalMaterialCost: 0,
				AcceptedAt:      ptr(ago(2 * time.Hour)),
				StartedAt:       ptr(ago(30 * time.Minute)),
			},
			Photos: []entity.OrderPhoto{
				{PhotoURL: "https://images.unsplash.com/photo-1585771724684-38269d6639fd?w=600", Caption: strPtr("Kebocoran di bawah wastafel"), DisplayOrder: 0},
				{PhotoURL: "https://images.unsplash.com/photo-1504328345606-18bbc8c9d7d1?w=600", Caption: strPtr("Pipa yang bocor"), DisplayOrder: 1},
			},
			Timeline: []entity.OrderTimeline{
				{Event: "order_created", Label: "Pesanan Dibuat", ActorType: strPtr("user")},
				{Event: "order_accepted", Label: "Pesanan Diterima", Description: strPtr("Tukang menerima, segera berangkat karena URGENT"), ActorType: strPtr("worker")},
				{Event: "worker_on_the_way", Label: "Tukang Berangkat", ActorType: strPtr("worker")},
				{Event: "worker_arrived", Label: "Tukang Tiba", ActorType: strPtr("worker")},
				{Event: "work_started", Label: "Pekerjaan Dimulai", Description: strPtr("Sedang mengidentifikasi sumber kebocoran"), ActorType: strPtr("worker")},
			},
			Chats: []entity.ChatMessage{
				{SenderID: userMap["andi.firmansyah@gmail.com"], SenderType: "user", Content: strPtr("Pak Hendra, URGENT! Pipa bocor di dapur, air terus mengalir. Bisa datang sekarang?"), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["hendra.wijaya@situkang.id"], SenderType: "worker", Content: strPtr("Siap Pak Andi, saya langsung berangkat sekarang. 15 menit lagi sampai."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["andi.firmansyah@gmail.com"], SenderType: "user", Content: strPtr("Makasih Pak, sementara sudah saya matikan stop kontak air utama."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["hendra.wijaya@situkang.id"], SenderType: "worker", Content: strPtr("Bagus Pak, itu tepat sekali. Saya sudah di depan ya."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["hendra.wijaya@situkang.id"], SenderType: "worker", Content: strPtr("Pak Andi, sumber bocornya ada di sambungan fleksibel wastafel yang sudah retak. Perlu ganti selang fleksibel + ring karet. Harga material sekitar Rp85.000. Boleh lanjut?"), MessageType: entity.MessageTypeText, IsRead: false},
			},
			Purchases: []entity.Purchase{
				{
					ItemName: "Selang Fleksibel Wastafel", Category: entity.PurchaseCategoryMaterial,
					Quantity: 1, Unit: "pcs", UnitPrice: 55000, TotalPrice: 55000,
					Reason:  strPtr("Selang fleksibel lama sudah retak dan jadi sumber kebocoran"),
					Status:  entity.PurchaseStatusPendingApproval, Confidence: float64Ptr(0.93),
					AIExplanation: strPtr("Item ini relevan dan sesuai dengan deskripsi pekerjaan perbaikan pipa bocor."),
				},
				{
					ItemName: "Ring Karet / Seal Pipa", Category: entity.PurchaseCategoryMaterial,
					Quantity: 3, Unit: "pcs", UnitPrice: 10000, TotalPrice: 30000,
					Reason:  strPtr("Ring karet perlu diganti agar sambungan tidak bocor lagi"),
					Status:  entity.PurchaseStatusPendingApproval, Confidence: float64Ptr(0.91),
				},
			},
		},

		// ── ORDER 4: Pending (Waiting Acceptance) ──────────────────────────────
		{
			Order: entity.Order{
				OrderNumber:     "HD-20260604-004",
				UserID:          userMap["siti.rahayu@gmail.com"],
				WorkerID:        userMap["rizki.permana@situkang.id"],
				ServiceID:       serviceMap["perbaikan-kusen"],
				CategoryID:      categoryMap["perbaikan-kusen"],
				Title:           "Pintu kamar susah ditutup, kusen bengkok",
				Description:     "Pintu kamar tidur susah ditutup sempurna. Kemungkinan kusen kayu sudah bengkok karena lembab. Perlu perbaikan atau penggantian kusen.",
				Status:          entity.OrderStatusPending,
				Urgency:         entity.OrderUrgencyNormal,
				LocationAddress: "Jl. Sudirman No. 55, Jakarta Pusat",
				LocationDetail:  strPtr("Kamar tidur anak, lantai 2"),
				LocationLat:     -6.2100,
				LocationLng:     106.8230,
				Notes:           strPtr("Pintu kayu jati, sudah 10 tahun"),
				BookingFee:      2000,
			},
			Photos: []entity.OrderPhoto{
				{PhotoURL: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600", Caption: strPtr("Kusen pintu yang bengkok"), DisplayOrder: 0},
			},
			Timeline: []entity.OrderTimeline{
				{Event: "order_created", Label: "Pesanan Dibuat", Description: strPtr("Menunggu konfirmasi dari tukang"), ActorType: strPtr("user")},
			},
			Chats: []entity.ChatMessage{
				{SenderID: userMap["siti.rahayu@gmail.com"], SenderType: "user", Content: strPtr("Pak Rizki, saya punya masalah pintu kamar yang susah ditutup. Bisa dicek dulu kondisinya?"), MessageType: entity.MessageTypeText, IsRead: false},
			},
		},

		// ── ORDER 5: Cancelled Order ───────────────────────────────────────────
		{
			Order: entity.Order{
				OrderNumber:     "HD-20260520-005",
				UserID:          userMap["dewi.kusuma@gmail.com"],
				WorkerID:        userMap["doni.prasetyo@situkang.id"],
				ServiceID:       serviceMap["pengecatan-rumah"],
				CategoryID:      categoryMap["pengecatan-rumah"],
				Title:           "Cat ulang 2 ruangan",
				Description:     "Ruang tamu dan ruang keluarga perlu dicat ulang. Warna sekarang sudah kusam.",
				Status:          entity.OrderStatusCancelled,
				Urgency:         entity.OrderUrgencyNormal,
				LocationAddress: "Jl. Kemang Raya No. 10, Jakarta Selatan",
				LocationLat:     -6.2600,
				LocationLng:     106.8150,
				BookingFee:      2000,
				CancellationReason:   strPtr("Ada keperluan mendadak keluarga, terpaksa reschedule"),
				CancellationCategory: strPtr("personal"),
				CancelledAt:          ptr(ago(3 * 24 * time.Hour)),
			},
			Timeline: []entity.OrderTimeline{
				{Event: "order_created", Label: "Pesanan Dibuat", ActorType: strPtr("user")},
				{Event: "order_accepted", Label: "Pesanan Diterima", ActorType: strPtr("worker")},
				{Event: "order_cancelled", Label: "Pesanan Dibatalkan", Description: strPtr("Dibatalkan oleh pelanggan karena keperluan mendadak"), ActorType: strPtr("user")},
			},
			Chats: []entity.ChatMessage{
				{SenderID: userMap["dewi.kusuma@gmail.com"], SenderType: "user", Content: strPtr("Pak Doni, maaf ya saya harus cancel ordernya. Ada keperluan keluarga mendadak."), MessageType: entity.MessageTypeText, IsRead: true},
				{SenderID: userMap["doni.prasetyo@situkang.id"], SenderType: "worker", Content: strPtr("Tidak apa-apa Bu Dewi, semoga keluarga baik-baik saja. Kalau sudah siap bisa order lagi ya."), MessageType: entity.MessageTypeText, IsRead: true},
			},
		},
	}

	for i := range seeds {
		s := &seeds[i]
		// Check duplicate order number
		var existing entity.Order
		if err := db.Where("order_number = ?", s.Order.OrderNumber).First(&existing).Error; err == nil {
			continue
		}

		// Create order
		if err := db.Create(&s.Order).Error; err != nil {
			return err
		}
		orderID := s.Order.ID

		// Create photos
		for j := range s.Photos {
			s.Photos[j].OrderID = orderID
			db.Create(&s.Photos[j])
		}

		// Create timeline
		for j := range s.Timeline {
			s.Timeline[j].OrderID = orderID
			if s.Timeline[j].ActorType != nil && *s.Timeline[j].ActorType == "user" {
				s.Timeline[j].ActorID = &s.Order.UserID
			} else if s.Timeline[j].ActorType != nil && *s.Timeline[j].ActorType == "worker" {
				s.Timeline[j].ActorID = &s.Order.WorkerID
			}
			db.Create(&s.Timeline[j])
		}

		// Create chats
		for j := range s.Chats {
			s.Chats[j].OrderID = orderID
			db.Create(&s.Chats[j])
		}

		// Create purchases
		for j := range s.Purchases {
			s.Purchases[j].OrderID = orderID
			s.Purchases[j].WorkerID = s.Order.WorkerID
			if s.Purchases[j].Status == entity.PurchaseStatusApproved {
				s.Purchases[j].ApprovedBy = &s.Order.UserID
				t := now.Add(-time.Duration(j) * time.Hour)
				s.Purchases[j].ApprovedAt = &t
			}
			db.Create(&s.Purchases[j])
		}

		// Create review (only for completed)
		if s.Review != nil && s.Order.Status == entity.OrderStatusCompleted {
			s.Review.OrderID = orderID
			s.Review.UserID = s.Order.UserID
			s.Review.WorkerID = s.Order.WorkerID
			if err := db.Where("order_id = ?", orderID).FirstOrCreate(s.Review).Error; err == nil {
				for _, tag := range s.ReviewTags {
					rt := entity.ReviewTag{ReviewID: s.Review.ID, Tag: tag}
					db.Where("review_id = ? AND tag = ?", s.Review.ID, tag).FirstOrCreate(&rt)
				}
			}
		}

		// Create invoice + payment for completed orders
		if s.Order.Status == entity.OrderStatusCompleted && s.Order.GrandTotal != nil {
			invoiceNum := "INV-" + s.Order.OrderNumber[3:]
			baseFee := 0
			if s.Order.BaseServiceFee != nil {
				baseFee = *s.Order.BaseServiceFee
			}
			invoice := entity.Invoice{
				OrderID:             orderID,
				InvoiceNumber:       invoiceNum,
				BaseServiceFee:      baseFee,
				TotalMaterialCost:   s.Order.TotalMaterialCost,
				BookingFee:          s.Order.BookingFee,
				GrandTotal:          *s.Order.GrandTotal,
				Currency:            "IDR",
				AIWorkSummary:       strPtr("Pekerjaan telah diselesaikan dengan baik sesuai deskripsi order."),
				AllPurchasesApproved: true,
			}
			var existInv entity.Invoice
			if err := db.Where("invoice_number = ?", invoiceNum).FirstOrCreate(&existInv, invoice).Error; err == nil {
				paidAt := *s.Order.CompletedAt
				payment := entity.Payment{
					OrderID:       orderID,
					InvoiceID:     existInv.ID,
					UserID:        s.Order.UserID,
					Amount:        *s.Order.GrandTotal,
					Currency:      "IDR",
					PaymentMethod: entity.PaymentMethodEWallet,
					PaymentStatus: entity.PaymentStatusPaid,
					TransactionRef: strPtr("TXN-" + s.Order.OrderNumber),
					PaidAt:        &paidAt,
				}
				var existPay entity.Payment
				db.Where("transaction_ref = ?", payment.TransactionRef).FirstOrCreate(&existPay, payment)
			}
		}
	}

	return nil
}
