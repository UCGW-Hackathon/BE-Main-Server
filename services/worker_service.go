package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"situkang/dto"
	"situkang/models/entity"
	http_error "situkang/models/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerService interface {
	GetProfile(ctx context.Context) (any, error)
	UpdateProfile(ctx context.Context, req dto.WorkerProfileUpdateRequest) (any, error)
	UpdateCoverPhoto(ctx context.Context, coverURL string) (any, error)
	SubmitVerification(ctx context.Context, idCardURL string, certificateURLs []string) (any, error)
	GetVerification(ctx context.Context) (any, error)
	GetHome(ctx context.Context) (any, error)
	UpdateAvailability(ctx context.Context, req dto.WorkerAvailabilityRequest) (any, error)
	ListIncomingOrders(ctx context.Context) (any, error)
	GetIncomingOrderDetail(ctx context.Context, orderID string) (any, error)
	AcceptOrder(ctx context.Context, orderID string, req dto.WorkerOrderAcceptRequest) (any, error)
	RejectOrder(ctx context.Context, orderID string, req dto.WorkerOrderRejectRequest) (any, error)
	ListOrders(ctx context.Context) (any, error)
	GetOrderDetail(ctx context.Context, orderID string) (any, error)
	UpdateOrderStatus(ctx context.Context, orderID string, req dto.WorkerOrderStatusRequest) (any, error)
	GenerateInvoice(ctx context.Context, orderID string, req dto.WorkerGenerateInvoiceRequest) (any, error)
	AddPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseCreateRequest) (any, error)
	AIProcessPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseAIProcessRequest) (any, error)
	ReceiptScanPurchase(ctx context.Context, orderID string, receiptURL string) (any, error)
	UpdatePurchase(ctx context.Context, orderID string, purchaseID string, req dto.WorkerPurchaseUpdateRequest) (any, error)
	DeletePurchase(ctx context.Context, orderID string, purchaseID string) error
	SubmitPurchase(ctx context.Context, orderID string, purchaseID string) (any, error)
	BulkSubmitPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseBulkSubmitRequest) (any, error)
	ClarifyPurchaseResponse(ctx context.Context, orderID string, purchaseID string, req dto.WorkerPurchaseClarifyResponseRequest) (any, error)
	ListChatMessages(ctx context.Context, orderID string) (any, error)
	SendChatMessage(ctx context.Context, orderID string, req dto.ChatSendRequest) (any, error)
	MarkChatRead(ctx context.Context, orderID string) error
	ListChats(ctx context.Context) (any, error)
	CreateCustomerRating(ctx context.Context, orderID string, req dto.CustomerRatingCreateRequest) (any, error)
	GetCustomerRating(ctx context.Context, orderID string) (any, error)
	GetHistory(ctx context.Context) (any, error)
	GetStatistics(ctx context.Context) (any, error)
	GetWallet(ctx context.Context) (any, error)
	ListWalletTransactions(ctx context.Context) (any, error)
	Withdraw(ctx context.Context, req dto.WalletWithdrawRequest) (any, error)
	UpdateLocation(ctx context.Context, req dto.WorkerLocationUpdateRequest) (any, error)
}

type workerService struct {
	db *gorm.DB
}

func NewWorkerService(db *gorm.DB) WorkerService {
	return &workerService{db: db}
}

func (s *workerService) GetProfile(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", workerID).Error; err != nil {
		return nil, err
	}
	profile, err := s.ensureWorkerProfile(ctx, workerID)
	if err != nil {
		return nil, err
	}
	services, _ := NewWorkerPublicService(s.db).GetServices(ctx, workerID.String())
	return workerProfileResponse(user, profile, services), nil
}

func (s *workerService) UpdateProfile(ctx context.Context, req dto.WorkerProfileUpdateRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.ensureWorkerProfileExists(ctx, workerID); err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userUpdates := map[string]any{}
		if req.FullName != nil {
			userUpdates["full_name"] = *req.FullName
		}
		if req.Phone != nil {
			userUpdates["phone"] = *req.Phone
		}
		if req.Address != nil {
			userUpdates["address"] = *req.Address
		}
		if len(userUpdates) > 0 {
			if err := tx.Model(&entity.User{}).Where("id = ?", workerID).Updates(userUpdates).Error; err != nil {
				return err
			}
		}
		profileUpdates := map[string]any{}
		if req.Specialization != nil {
			profileUpdates["specialization"] = *req.Specialization
		}
		if req.Bio != nil {
			profileUpdates["bio"] = *req.Bio
		}
		if req.BasePrice != nil {
			profileUpdates["base_price"] = *req.BasePrice
		}
		if req.PriceUnit != nil {
			profileUpdates["price_unit"] = *req.PriceUnit
		}
		if len(profileUpdates) > 0 {
			if err := tx.Model(&entity.WorkerProfile{}).Where("user_id = ?", workerID).Updates(profileUpdates).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return s.GetProfile(ctx)
}

func (s *workerService) UpdateCoverPhoto(ctx context.Context, coverURL string) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.ensureWorkerProfileExists(ctx, workerID); err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.WorkerProfile{}).Where("user_id = ?", workerID).Update("cover_photo_url", coverURL).Error; err != nil {
		return nil, err
	}
	return map[string]any{"cover_photo_url": coverURL}, nil
}

func (s *workerService) SubmitVerification(ctx context.Context, idCardURL string, certificateURLs []string) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.ensureWorkerProfileExists(ctx, workerID); err != nil {
		return nil, err
	}
	certs := entity.JSONB("[]")
	if certificateURLs != nil {
		raw, _ := json.Marshal(certificateURLs)
		certs = raw
	}
	if err := s.db.WithContext(ctx).Model(&entity.WorkerProfile{}).Where("user_id = ?", workerID).Updates(map[string]any{
		"verification_status": entity.VerificationStatusPending,
		"id_card_url":         idCardURL,
		"certificate_urls":    certs,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"verification_status": entity.VerificationStatusPending, "id_card_url": idCardURL, "certificate_urls": certs}, nil
}

func (s *workerService) GetVerification(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	profile, err := s.ensureWorkerProfile(ctx, workerID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"verification_status": profile.VerificationStatus,
		"id_card_url":         profile.IDCardURL,
		"certificate_urls":    profile.CertificateURLs,
		"verified_at":         profile.VerifiedAt,
	}, nil
}

func (s *workerService) GetHome(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", workerID).Error; err != nil {
		return nil, err
	}
	profile, _ := s.ensureWorkerProfile(ctx, workerID)
	var wallet entity.WorkerWallet
	_ = s.db.WithContext(ctx).First(&wallet, "worker_id = ?", workerID).Error
	var incomingCount, activeCount int64
	_ = s.db.WithContext(ctx).Model(&entity.Order{}).Where("worker_id = ? AND status = ?", workerID, entity.OrderStatusPending).Count(&incomingCount).Error
	_ = s.db.WithContext(ctx).Model(&entity.Order{}).Where("worker_id = ? AND status IN ?", workerID, []entity.OrderStatus{entity.OrderStatusAccepted, entity.OrderStatusOnTheWay, entity.OrderStatusArrived, entity.OrderStatusInProgress}).Count(&activeCount).Error
	return map[string]any{
		"worker_summary": map[string]any{
			"full_name":      user.FullName,
			"avatar_url":     user.AvatarURL,
			"is_available":   profile.IsAvailable,
			"rating":         profile.RatingAvg,
			"completed_jobs": profile.CompletedJobs,
			"balance":        wallet.Balance,
		},
		"incoming_orders_count": incomingCount,
		"active_orders_count":   activeCount,
	}, nil
}

func (s *workerService) UpdateAvailability(ctx context.Context, req dto.WorkerAvailabilityRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.ensureWorkerProfileExists(ctx, workerID); err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.WorkerProfile{}).Where("user_id = ?", workerID).Update("is_available", req.IsAvailable).Error; err != nil {
		return nil, err
	}
	return map[string]any{"is_available": req.IsAvailable}, nil
}

func (s *workerService) ListIncomingOrders(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var orders []entity.Order
	if err := s.db.WithContext(ctx).Where("worker_id = ? AND status = ?", workerID, entity.OrderStatusPending).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return NewOrderService(s.db).(*orderService).orderList(ctx, orders), nil
}

func (s *workerService) GetIncomingOrderDetail(ctx context.Context, orderID string) (any, error) {
	return s.GetOrderDetail(ctx, orderID)
}

func (s *workerService) AcceptOrder(ctx context.Context, orderID string, req dto.WorkerOrderAcceptRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	if order.Status != entity.OrderStatusPending {
		return nil, http_error.INVALID_STATUS_TRANSITION
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(map[string]any{"status": entity.OrderStatusAccepted, "accepted_at": &now}).Error; err != nil {
			return err
		}
		return createTimeline(tx, order.ID, "accepted", "Pesanan diterima", req.Note, &workerID, "worker")
	}); err != nil {
		return nil, err
	}
	return map[string]any{"order_id": order.ID.String(), "status": entity.OrderStatusAccepted, "accepted_at": now}, nil
}

func (s *workerService) RejectOrder(ctx context.Context, orderID string, req dto.WorkerOrderRejectRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	if order.Status != entity.OrderStatusPending {
		return nil, http_error.INVALID_STATUS_TRANSITION
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&entity.Order{}).Where("id = ?", order.ID).Updates(map[string]any{
		"status":                entity.OrderStatusRejected,
		"cancellation_reason":   req.Reason,
		"cancellation_category": req.ReasonCategory,
		"cancelled_by":          workerID,
		"cancelled_at":          &now,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"order_id": order.ID.String(), "status": entity.OrderStatusRejected}, nil
}

func (s *workerService) ListOrders(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var orders []entity.Order
	if err := s.db.WithContext(ctx).Where("worker_id = ?", workerID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return NewOrderService(s.db).(*orderService).orderList(ctx, orders), nil
}

func (s *workerService) GetOrderDetail(ctx context.Context, orderID string) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	return NewOrderService(s.db).(*orderService).orderDetail(ctx, order)
}

func (s *workerService) UpdateOrderStatus(ctx context.Context, orderID string, req dto.WorkerOrderStatusRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	newStatus := entity.OrderStatus(req.Status)
	if !isKnownOrderStatus(newStatus) {
		return nil, http_error.INVALID_STATUS_TRANSITION
	}
	updates := map[string]any{"status": newStatus}
	now := time.Now()
	if newStatus == entity.OrderStatusInProgress && order.StartedAt == nil {
		updates["started_at"] = &now
	}
	if newStatus == entity.OrderStatusCompleted && order.CompletedAt == nil {
		updates["completed_at"] = &now
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
			return err
		}
		return createTimeline(tx, order.ID, string(newStatus), "Status pesanan diperbarui", req.Note, &workerID, "worker")
	}); err != nil {
		return nil, err
	}
	return map[string]any{"order_id": order.ID.String(), "status": newStatus, "updated_at": now}, nil
}

func (s *workerService) GenerateInvoice(ctx context.Context, orderID string, req dto.WorkerGenerateInvoiceRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	var purchases []entity.Purchase
	_ = s.db.WithContext(ctx).Where("order_id = ? AND status = ?", order.ID, entity.PurchaseStatusApproved).Find(&purchases).Error
	materialTotal := 0
	for _, purchase := range purchases {
		materialTotal += purchase.TotalPrice
	}
	grandTotal := req.BaseServiceFee + materialTotal + order.BookingFee
	invoice := entity.Invoice{
		OrderID:              order.ID,
		InvoiceNumber:        fmt.Sprintf("INV-%s", order.OrderNumber),
		BaseServiceFee:       req.BaseServiceFee,
		TotalMaterialCost:    materialTotal,
		BookingFee:           order.BookingFee,
		GrandTotal:           grandTotal,
		WorkerNotes:          req.WorkerNotes,
		AllPurchasesApproved: true,
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", order.ID).Assign(invoice).FirstOrCreate(&invoice).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(map[string]any{
			"base_service_fee":    req.BaseServiceFee,
			"total_material_cost": materialTotal,
			"grand_total":         grandTotal,
		}).Error
	}); err != nil {
		return nil, err
	}
	return map[string]any{"invoice_id": invoice.ID.String(), "order_id": order.ID.String(), "invoice_number": invoice.InvoiceNumber, "grand_total": invoice.GrandTotal}, nil
}

func (s *workerService) AddPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseCreateRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	purchase := entity.Purchase{
		OrderID:         order.ID,
		WorkerID:        workerID,
		ItemName:        req.ItemName,
		Category:        entity.PurchaseCategory(req.Category),
		Quantity:        req.Quantity,
		Unit:            req.Unit,
		UnitPrice:       req.UnitPrice,
		TotalPrice:      req.TotalPrice,
		Reason:          req.Reason,
		ReceiptPhotoURL: req.ReceiptPhotoURL,
		Status:          entity.PurchaseStatusDraft,
	}
	if purchase.TotalPrice == 0 {
		purchase.TotalPrice = int(purchase.Quantity * float64(purchase.UnitPrice))
	}
	if err := s.db.WithContext(ctx).Create(&purchase).Error; err != nil {
		return nil, err
	}
	return purchaseResponse(purchase), nil
}

func (s *workerService) AIProcessPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseAIProcessRequest) (any, error) {
	return map[string]any{
		"order_id":        orderID,
		"items":           []any{},
		"summary":         req.RawInput,
		"risk_flags":      []any{},
		"approval_status": entity.PurchaseStatusDraft,
	}, nil
}

func (s *workerService) ReceiptScanPurchase(ctx context.Context, orderID string, receiptURL string) (any, error) {
	return map[string]any{"order_id": orderID, "receipt_photo_url": receiptURL, "items": []any{}, "risk_flags": []any{}}, nil
}

func (s *workerService) UpdatePurchase(ctx context.Context, orderID string, purchaseID string, req dto.WorkerPurchaseUpdateRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	purchase, err := NewOrderService(s.db).(*orderService).getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	if purchase.Status == entity.PurchaseStatusApproved {
		return nil, http_error.INVALID_STATUS_TRANSITION
	}
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Updates(map[string]any{
		"item_name":   req.ItemName,
		"category":    entity.PurchaseCategory(req.Category),
		"quantity":    req.Quantity,
		"unit":        req.Unit,
		"unit_price":  req.UnitPrice,
		"total_price": req.TotalPrice,
		"reason":      req.Reason,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "order_id": order.ID.String()}, nil
}

func (s *workerService) DeletePurchase(ctx context.Context, orderID string, purchaseID string) error {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return err
	}
	purchase, err := NewOrderService(s.db).(*orderService).getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return err
	}
	if purchase.Status == entity.PurchaseStatusApproved {
		return http_error.INVALID_STATUS_TRANSITION
	}
	return s.db.WithContext(ctx).Delete(&purchase).Error
}

func (s *workerService) SubmitPurchase(ctx context.Context, orderID string, purchaseID string) (any, error) {
	return s.updatePurchaseStatus(ctx, orderID, purchaseID, entity.PurchaseStatusPendingApproval)
}

func (s *workerService) BulkSubmitPurchase(ctx context.Context, orderID string, req dto.WorkerPurchaseBulkSubmitRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	ids, err := parseUUIDList(req.PurchaseIDs)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("order_id = ? AND id IN ?", order.ID, ids).Update("status", entity.PurchaseStatusPendingApproval).Error; err != nil {
		return nil, err
	}
	return map[string]any{"submitted_count": len(ids), "submitted_ids": req.PurchaseIDs}, nil
}

func (s *workerService) ClarifyPurchaseResponse(ctx context.Context, orderID string, purchaseID string, req dto.WorkerPurchaseClarifyResponseRequest) (any, error) {
	updates := map[string]any{
		"status":                 entity.PurchaseStatusPendingApproval,
		"needs_clarification":    false,
		"clarification_response": req.Response,
	}
	if req.UpdatedItemName != nil {
		updates["item_name"] = *req.UpdatedItemName
	}
	if req.UpdatedReason != nil {
		updates["reason"] = *req.UpdatedReason
	}
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	purchase, err := NewOrderService(s.db).(*orderService).getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "status": entity.PurchaseStatusPendingApproval}, nil
}

func (s *workerService) ListChatMessages(ctx context.Context, orderID string) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	var messages []entity.ChatMessage
	if err := s.db.WithContext(ctx).Where("order_id = ?", order.ID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	data := make([]map[string]any, 0, len(messages))
	for _, message := range messages {
		data = append(data, chatResponse(message))
	}
	return map[string]any{"order_id": order.ID.String(), "messages": data, "has_more": false}, nil
}

func (s *workerService) SendChatMessage(ctx context.Context, orderID string, req dto.ChatSendRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	message := entity.ChatMessage{OrderID: order.ID, SenderID: workerID, SenderType: "worker", Content: req.Content, MessageType: entity.MessageType(req.MessageType)}
	if err := s.db.WithContext(ctx).Create(&message).Error; err != nil {
		return nil, err
	}
	return chatResponse(message), nil
}

func (s *workerService) MarkChatRead(ctx context.Context, orderID string) error {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return err
	}
	now := time.Now()
	return s.db.WithContext(ctx).Model(&entity.ChatMessage{}).Where("order_id = ? AND sender_id <> ? AND is_read = FALSE", order.ID, workerID).Updates(map[string]any{"is_read": true, "read_at": &now}).Error
}

func (s *workerService) ListChats(ctx context.Context) (any, error) {
	return s.ListOrders(ctx)
}

func (s *workerService) CreateCustomerRating(ctx context.Context, orderID string, req dto.CustomerRatingCreateRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.Rating < 1 || req.Rating > 5 {
		return nil, http_error.VALIDATION_ERROR
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	review := entity.Review{OrderID: order.ID, UserID: order.UserID, WorkerID: workerID, ReviewType: "customer", Rating: int16(req.Rating), Comment: req.Comment}
	if err := NewOrderService(s.db).(*orderService).createReviewWithTags(ctx, review, req.Tags); err != nil {
		return nil, err
	}
	return map[string]any{"customer_review_id": review.ID.String(), "order_id": order.ID.String()}, nil
}

func (s *workerService) GetCustomerRating(ctx context.Context, orderID string) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	var review entity.Review
	if err := s.db.WithContext(ctx).First(&review, "order_id = ? AND review_type = ?", order.ID, "customer").Error; err != nil {
		return nil, err
	}
	return reviewResponse(review), nil
}

func (s *workerService) GetHistory(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var orders []entity.Order
	if err := s.db.WithContext(ctx).Where("worker_id = ? AND status IN ?", workerID, []entity.OrderStatus{entity.OrderStatusCompleted, entity.OrderStatusCancelled, entity.OrderStatusRejected}).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return map[string]any{"orders": NewOrderService(s.db).(*orderService).orderList(ctx, orders)}, nil
}

func (s *workerService) GetStatistics(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var completed int64
	var earnings int64
	_ = s.db.WithContext(ctx).Model(&entity.Order{}).Where("worker_id = ? AND status = ?", workerID, entity.OrderStatusCompleted).Count(&completed).Error
	_ = s.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Joins("JOIN worker_wallets ww ON ww.id = wallet_transactions.wallet_id").
		Where("ww.worker_id = ? AND wallet_transactions.type = ? AND wallet_transactions.status = ?", workerID, entity.WalletTxTypeEarning, entity.WalletTxStatusCompleted).
		Select("COALESCE(SUM(wallet_transactions.amount), 0)").Scan(&earnings).Error
	return map[string]any{"completed_jobs": completed, "total_earnings": earnings}, nil
}

func (s *workerService) GetWallet(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	wallet, err := s.ensureWallet(ctx, workerID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"wallet_id": wallet.ID.String(), "balance": wallet.Balance, "total_earnings": wallet.TotalEarnings, "total_withdrawn": wallet.TotalWithdrawn, "pending_earnings": wallet.PendingEarnings, "is_active": wallet.IsActive, "currency": "IDR", "updated_at": wallet.UpdatedAt}, nil
}

func (s *workerService) ListWalletTransactions(ctx context.Context) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	wallet, err := s.ensureWallet(ctx, workerID)
	if err != nil {
		return nil, err
	}
	var transactions []entity.WalletTransaction
	if err := s.db.WithContext(ctx).Where("wallet_id = ?", wallet.ID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	data := make([]map[string]any, 0, len(transactions))
	for _, tx := range transactions {
		data = append(data, map[string]any{"transaction_id": tx.ID.String(), "type": tx.Type, "amount": tx.Amount, "balance_before": tx.BalanceBefore, "balance_after": tx.BalanceAfter, "description": tx.Description, "order_id": tx.OrderID, "status": tx.Status, "completed_at": tx.CompletedAt, "created_at": tx.CreatedAt})
	}
	return data, nil
}

func (s *workerService) Withdraw(ctx context.Context, req dto.WalletWithdrawRequest) (any, error) {
	if req.Amount < 50000 {
		return nil, http_error.VALIDATION_ERROR
	}
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	wallet, err := s.ensureWallet(ctx, workerID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < int64(req.Amount) {
		return nil, http_error.INSUFFICIENT_BALANCE
	}
	description := fmt.Sprintf("Penarikan ke rekening %s ****%s", req.BankName, last4(req.AccountNumber))
	transaction := entity.WalletTransaction{
		WalletID:      wallet.ID,
		Type:          entity.WalletTxTypeWithdrawal,
		Amount:        req.Amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance - int64(req.Amount),
		Description:   &description,
		Status:        entity.WalletTxStatusPending,
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return tx.Model(&entity.WorkerWallet{}).Where("id = ?", wallet.ID).Updates(map[string]any{
			"balance":         transaction.BalanceAfter,
			"total_withdrawn": wallet.TotalWithdrawn + int64(req.Amount),
		}).Error
	}); err != nil {
		return nil, err
	}
	return map[string]any{"transaction_id": transaction.ID.String(), "type": transaction.Type, "amount": transaction.Amount, "balance_before": transaction.BalanceBefore, "balance_after": transaction.BalanceAfter, "status": transaction.Status, "created_at": transaction.CreatedAt}, nil
}

func (s *workerService) UpdateLocation(ctx context.Context, req dto.WorkerLocationUpdateRequest) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", workerID).Updates(map[string]any{"latitude": req.Latitude, "longitude": req.Longitude}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"latitude": req.Latitude, "longitude": req.Longitude, "eta_minutes": nil, "distance_remaining_km": nil, "updated_at": time.Now()}, nil
}

func (s *workerService) getOrderForWorker(ctx context.Context, orderID string, workerID uuid.UUID) (entity.Order, error) {
	id, err := parseUUID(orderID)
	if err != nil {
		return entity.Order{}, err
	}
	var order entity.Order
	err = s.db.WithContext(ctx).First(&order, "id = ? AND worker_id = ?", id, workerID).Error
	return order, err
}

func (s *workerService) updatePurchaseStatus(ctx context.Context, orderID string, purchaseID string, status entity.PurchaseStatus) (any, error) {
	workerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForWorker(ctx, orderID, workerID)
	if err != nil {
		return nil, err
	}
	purchase, err := NewOrderService(s.db).(*orderService).getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Update("status", status).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "status": status}, nil
}

func (s *workerService) ensureWorkerProfile(ctx context.Context, workerID uuid.UUID) (entity.WorkerProfile, error) {
	var profile entity.WorkerProfile
	err := s.db.WithContext(ctx).First(&profile, "user_id = ?", workerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = entity.WorkerProfile{UserID: workerID, VerificationStatus: entity.VerificationStatusUnverified, CertificateURLs: entity.JSONB("[]"), IsAvailable: true}
		err = s.db.WithContext(ctx).Create(&profile).Error
	}
	return profile, err
}

func (s *workerService) ensureWorkerProfileExists(ctx context.Context, workerID uuid.UUID) error {
	_, err := s.ensureWorkerProfile(ctx, workerID)
	return err
}

func (s *workerService) ensureWallet(ctx context.Context, workerID uuid.UUID) (entity.WorkerWallet, error) {
	var wallet entity.WorkerWallet
	err := s.db.WithContext(ctx).First(&wallet, "worker_id = ?", workerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		wallet = entity.WorkerWallet{WorkerID: workerID, IsActive: true}
		err = s.db.WithContext(ctx).Create(&wallet).Error
	}
	return wallet, err
}

func workerProfileResponse(user entity.User, profile entity.WorkerProfile, services any) map[string]any {
	return map[string]any{
		"worker_id":           user.ID.String(),
		"full_name":           user.FullName,
		"email":               user.Email,
		"phone":               user.Phone,
		"avatar_url":          user.AvatarURL,
		"address":             user.Address,
		"latitude":            user.Latitude,
		"longitude":           user.Longitude,
		"specialization":      profile.Specialization,
		"bio":                 profile.Bio,
		"cover_photo_url":     profile.CoverPhotoURL,
		"verification_status": profile.VerificationStatus,
		"base_price":          profile.BasePrice,
		"price_unit":          profile.PriceUnit,
		"booking_fee":         profile.BookingFee,
		"rating":              profile.RatingAvg,
		"total_reviews":       profile.TotalReviews,
		"completed_jobs":      profile.CompletedJobs,
		"is_available":        profile.IsAvailable,
		"services":            services,
	}
}

func isKnownOrderStatus(status entity.OrderStatus) bool {
	switch status {
	case entity.OrderStatusPending, entity.OrderStatusAccepted, entity.OrderStatusOnTheWay, entity.OrderStatusArrived, entity.OrderStatusInProgress, entity.OrderStatusWorkPaused, entity.OrderStatusCompleted, entity.OrderStatusCancelled, entity.OrderStatusRejected:
		return true
	default:
		return false
	}
}

func last4(value string) string {
	if len(value) <= 4 {
		return value
	}
	return value[len(value)-4:]
}
