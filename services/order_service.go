package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"whatsapp-backend/dto"
	"whatsapp-backend/models/entity"
	http_error "whatsapp-backend/models/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (any, error)
	ListOrders(ctx context.Context) (any, error)
	GetOrderDetail(ctx context.Context, orderID string) (any, error)
	CancelOrder(ctx context.Context, orderID string, req dto.OrderCancelRequest) (any, error)
	GetTracking(ctx context.Context, orderID string) (any, error)
	GetTrackingLocation(ctx context.Context, orderID string) (any, error)
	ListPurchases(ctx context.Context, orderID string) (any, error)
	GetPurchaseDetail(ctx context.Context, orderID string, purchaseID string) (any, error)
	ApprovePurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseApproveRequest) (any, error)
	RejectPurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseRejectRequest) (any, error)
	ClarifyPurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseClarifyRequest) (any, error)
	BulkApprovePurchases(ctx context.Context, orderID string, req dto.PurchaseBulkApproveRequest) (any, error)
	ListChatMessages(ctx context.Context, orderID string) (any, error)
	SendChatMessage(ctx context.Context, orderID string, req dto.ChatSendRequest) (any, error)
	MarkChatRead(ctx context.Context, orderID string) error
	ListChats(ctx context.Context) (any, error)
	CreateRating(ctx context.Context, orderID string, req dto.RatingCreateRequest) (any, error)
	GetRating(ctx context.Context, orderID string) (any, error)
	GetInvoice(ctx context.Context, orderID string) (any, error)
	CreatePayment(ctx context.Context, orderID string, req dto.PaymentCreateRequest) (any, error)
	DownloadInvoicePDF(ctx context.Context, orderID string) ([]byte, error)
	SandboxCallback(ctx context.Context, req dto.SandboxPaymentCallbackRequest) (any, error)
	GetPaymentDetails(ctx context.Context, paymentID string) (any, error)
}

type orderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderService {
	return &orderService{db: db}
}

func (s *orderService) CreateOrder(ctx context.Context, req dto.OrderCreateRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	workerID, err := parseUUID(req.WorkerID)
	if err != nil {
		return nil, err
	}
	serviceID, err := parseUUID(req.ServiceID)
	if err != nil {
		return nil, err
	}

	var service entity.Service
	if err := s.db.WithContext(ctx).First(&service, "id = ? AND is_active = TRUE", serviceID).Error; err != nil {
		return nil, err
	}
	var workerProfile entity.WorkerProfile
	if err := s.db.WithContext(ctx).First(&workerProfile, "user_id = ?", workerID).Error; err != nil {
		return nil, err
	}

	urgency := entity.OrderUrgencyNormal
	if req.Urgency != nil && *req.Urgency == string(entity.OrderUrgencyUrgent) {
		urgency = entity.OrderUrgencyUrgent
	}
	preferredDate, err := parseOptionalDate(req.PreferredDate)
	if err != nil {
		return nil, err
	}
	preferredStart, err := parseOptionalClock(req.PreferredTimeStart)
	if err != nil {
		return nil, err
	}
	preferredEnd, err := parseOptionalClock(req.PreferredTimeEnd)
	if err != nil {
		return nil, err
	}

	baseServiceFee := service.BasePrice
	if baseServiceFee == nil {
		baseServiceFee = workerProfile.BasePrice
	}

	order := entity.Order{
		OrderNumber:        newOrderNumber(),
		UserID:             userID,
		WorkerID:           workerID,
		ServiceID:          service.ID,
		CategoryID:         service.CategoryID,
		Title:              req.Title,
		Description:        req.Description,
		Status:             entity.OrderStatusPending,
		Urgency:            urgency,
		LocationAddress:    req.Location.Address,
		LocationDetail:     req.Location.AddressDetail,
		LocationLat:        req.Location.Latitude,
		LocationLng:        req.Location.Longitude,
		PreferredDate:      preferredDate,
		PreferredTimeStart: preferredStart,
		PreferredTimeEnd:   preferredEnd,
		Notes:              req.Notes,
		BookingFee:         workerProfile.BookingFee,
		BaseServiceFee:     baseServiceFee,
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for idx, photo := range req.Photos {
			orderPhoto := entity.OrderPhoto{
				OrderID:      order.ID,
				PhotoURL:     photo,
				DisplayOrder: idx + 1,
			}
			if err := tx.Create(&orderPhoto).Error; err != nil {
				return err
			}
		}
		return createTimeline(tx, order.ID, "created", "Pesanan dibuat", nil, &userID, "user")
	}); err != nil {
		return nil, err
	}

	return map[string]any{
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"status":       order.Status,
		"booking_fee":  order.BookingFee,
		"created_at":   order.CreatedAt,
	}, nil
}

func (s *orderService) ListOrders(ctx context.Context) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var orders []entity.Order
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return s.orderList(ctx, orders), nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	return s.orderDetail(ctx, order)
}

func (s *orderService) CancelOrder(ctx context.Context, orderID string, req dto.OrderCancelRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	if order.Status == entity.OrderStatusArrived ||
		order.Status == entity.OrderStatusInProgress ||
		order.Status == entity.OrderStatusCompleted {
		return nil, http_error.INVALID_STATUS_TRANSITION
	}

	now := time.Now()
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(map[string]any{
			"status":                entity.OrderStatusCancelled,
			"cancellation_reason":   req.Reason,
			"cancellation_category": req.ReasonCategory,
			"cancelled_by":          userID,
			"cancelled_at":          &now,
		}).Error; err != nil {
			return err
		}
		return createTimeline(tx, order.ID, "cancelled", "Pesanan dibatalkan", &req.Reason, &userID, "user")
	}); err != nil {
		return nil, err
	}

	return map[string]any{
		"order_id":     order.ID.String(),
		"order_number": order.OrderNumber,
		"status":       entity.OrderStatusCancelled,
		"cancelled_at": now,
	}, nil
}

func (s *orderService) GetTracking(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	var timeline []entity.OrderTimeline
	if err := s.db.WithContext(ctx).Where("order_id = ?", order.ID).Order("created_at ASC").Find(&timeline).Error; err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(timeline))
	for _, event := range timeline {
		items = append(items, map[string]any{
			"event":       event.Event,
			"label":       event.Label,
			"description": event.Description,
			"metadata":    event.Metadata,
			"created_at":  event.CreatedAt,
		})
	}
	location, _ := s.GetTrackingLocation(ctx, orderID)
	return map[string]any{
		"order_id":        order.ID.String(),
		"order_number":    order.OrderNumber,
		"status":          order.Status,
		"worker_location": location,
		"timeline":        items,
	}, nil
}

func (s *orderService) GetTrackingLocation(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	var worker entity.User
	if err := s.db.WithContext(ctx).First(&worker, "id = ?", order.WorkerID).Error; err != nil {
		return nil, err
	}
	distance := distanceKm(order.LocationLat, order.LocationLng, worker.Latitude, worker.Longitude)
	return map[string]any{
		"worker_id":             worker.ID.String(),
		"latitude":              worker.Latitude,
		"longitude":             worker.Longitude,
		"distance_remaining_km": distance,
		"eta_minutes":           estimateETA(distance),
		"updated_at":            worker.UpdatedAt,
	}, nil
}

func (s *orderService) ListPurchases(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}

	var purchases []entity.Purchase
	if err := s.db.WithContext(ctx).Where("order_id = ?", order.ID).Order("created_at DESC").Find(&purchases).Error; err != nil {
		return nil, err
	}
	totalApproved := 0
	items := make([]map[string]any, 0, len(purchases))
	for _, purchase := range purchases {
		if purchase.Status == entity.PurchaseStatusApproved {
			totalApproved += purchase.TotalPrice
		}
		items = append(items, purchaseResponse(purchase))
	}
	return map[string]any{
		"order_id":  order.ID.String(),
		"purchases": items,
		"summary": map[string]any{
			"count":               len(items),
			"total_approved_cost": totalApproved,
		},
	}, nil
}

func (s *orderService) GetPurchaseDetail(ctx context.Context, orderID string, purchaseID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	purchase, err := s.getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	data := purchaseResponse(purchase)
	var flags []entity.PurchaseRiskFlag
	_ = s.db.WithContext(ctx).Where("purchase_id = ?", purchase.ID).Find(&flags).Error
	data["risk_flags"] = flags
	return data, nil
}

func (s *orderService) ApprovePurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseApproveRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	purchase, err := s.getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Updates(map[string]any{
		"status":      entity.PurchaseStatusApproved,
		"approved_by": userID,
		"approved_at": &now,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "status": entity.PurchaseStatusApproved, "approved_at": now}, nil
}

func (s *orderService) RejectPurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseRejectRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	purchase, err := s.getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Updates(map[string]any{
		"status":           entity.PurchaseStatusRejected,
		"rejected_by":      userID,
		"rejected_at":      &now,
		"rejection_reason": req.Reason,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "status": entity.PurchaseStatusRejected, "rejected_at": now}, nil
}

func (s *orderService) ClarifyPurchase(ctx context.Context, orderID string, purchaseID string, req dto.PurchaseClarifyRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	purchase, err := s.getPurchase(ctx, order.ID, purchaseID)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).Where("id = ?", purchase.ID).Updates(map[string]any{
		"status":                 entity.PurchaseStatusNeedsClarification,
		"needs_clarification":    true,
		"clarification_question": req.Question,
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"purchase_id": purchase.ID.String(), "status": entity.PurchaseStatusNeedsClarification}, nil
}

func (s *orderService) BulkApprovePurchases(ctx context.Context, orderID string, req dto.PurchaseBulkApproveRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	ids, err := parseUUIDList(req.PurchaseIDs)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&entity.Purchase{}).
		Where("order_id = ? AND id IN ?", order.ID, ids).
		Updates(map[string]any{"status": entity.PurchaseStatusApproved, "approved_by": userID, "approved_at": &now}).Error; err != nil {
		return nil, err
	}
	return map[string]any{"approved_count": len(ids), "approved_ids": req.PurchaseIDs}, nil
}

func (s *orderService) ListChatMessages(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
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

func (s *orderService) SendChatMessage(ctx context.Context, orderID string, req dto.ChatSendRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	role, _ := currentRole(ctx)
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	message := entity.ChatMessage{
		OrderID:     order.ID,
		SenderID:    userID,
		SenderType:  role,
		Content:     req.Content,
		MessageType: entity.MessageType(req.MessageType),
	}
	if err := s.db.WithContext(ctx).Create(&message).Error; err != nil {
		return nil, err
	}
	return chatResponse(message), nil
}

func (s *orderService) MarkChatRead(ctx context.Context, orderID string) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return err
	}
	now := time.Now()
	return s.db.WithContext(ctx).Model(&entity.ChatMessage{}).
		Where("order_id = ? AND sender_id <> ? AND is_read = FALSE", order.ID, userID).
		Updates(map[string]any{"is_read": true, "read_at": &now}).Error
}

func (s *orderService) ListChats(ctx context.Context) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var orders []entity.Order
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return s.orderList(ctx, orders), nil
}

func (s *orderService) CreateRating(ctx context.Context, orderID string, req dto.RatingCreateRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.Rating < 1 || req.Rating > 5 {
		return nil, http_error.VALIDATION_ERROR
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	review := entity.Review{
		OrderID:    order.ID,
		UserID:     userID,
		WorkerID:   order.WorkerID,
		ReviewType: "worker",
		Rating:     int16(req.Rating),
		Comment:    req.Comment,
	}
	if err := s.createReviewWithTags(ctx, review, req.Tags); err != nil {
		return nil, err
	}
	return map[string]any{"review_id": review.ID.String(), "order_id": order.ID.String()}, nil
}

func (s *orderService) GetRating(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	var review entity.Review
	if err := s.db.WithContext(ctx).First(&review, "order_id = ? AND review_type = ?", order.ID, "worker").Error; err != nil {
		return nil, err
	}
	return reviewResponse(review), nil
}

func (s *orderService) GetInvoice(ctx context.Context, orderID string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	return s.invoiceResponse(ctx, order.ID)
}

func (s *orderService) CreatePayment(ctx context.Context, orderID string, req dto.PaymentCreateRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	var invoice entity.Invoice
	if err := s.db.WithContext(ctx).First(&invoice, "order_id = ?", order.ID).Error; err != nil {
		return nil, err
	}
	paymentRef := fmt.Sprintf("PAY-REF-%s", uuid.NewString()[:8])
	payment := entity.Payment{
		OrderID:         order.ID,
		InvoiceID:       invoice.ID,
		UserID:          userID,
		Amount:          invoice.GrandTotal,
		PaymentMethod:   entity.PaymentMethod(req.PaymentMethod),
		PaymentStatus:   entity.PaymentStatusPending,
		PaymentProofURL: req.PaymentProofURL,
		TransactionRef:  &paymentRef,
	}
	if err := s.db.WithContext(ctx).Create(&payment).Error; err != nil {
		return nil, err
	}
	return map[string]any{
		"payment_id":     payment.ID.String(),
		"order_id":       order.ID.String(),
		"invoice_id":     invoice.ID.String(),
		"amount":         payment.Amount,
		"payment_status": payment.PaymentStatus,
		"token":          paymentRef,
		"redirect_url":   fmt.Sprintf("http://localhost:8080/v1/payments/sandbox-checkout?payment_id=%s", payment.ID.String()),
		"created_at":     payment.CreatedAt,
	}, nil
}

func (s *orderService) SandboxCallback(ctx context.Context, req dto.SandboxPaymentCallbackRequest) (any, error) {
	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return nil, http_error.BAD_REQUEST_ERROR
	}

	var payment entity.Payment
	if err := s.db.WithContext(ctx).First(&payment, "id = ?", paymentID).Error; err != nil {
		return nil, err
	}

	if payment.PaymentStatus == entity.PaymentStatusPaid {
		return map[string]any{
			"payment_id":     payment.ID.String(),
			"order_id":       payment.OrderID.String(),
			"payment_status": payment.PaymentStatus,
			"paid_at":        payment.PaidAt,
		}, nil
	}

	now := time.Now()
	paymentStatus := entity.PaymentStatusPaid
	if req.Status == "failed" {
		paymentStatus = entity.PaymentStatusRefunded
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payment.PaymentStatus = paymentStatus
		payment.PaidAt = &now
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		if paymentStatus == entity.PaymentStatusPaid {
			if err := tx.Model(&entity.Order{}).Where("id = ?", payment.OrderID).Updates(map[string]any{
				"status":       entity.OrderStatusCompleted,
				"completed_at": &now,
			}).Error; err != nil {
				return err
			}

			desc := "Pembayaran berhasil diverifikasi secara instan via Payment Gateway Sandbox"
			timeline := entity.OrderTimeline{
				OrderID:     payment.OrderID,
				Event:       "paid",
				Label:       "Pembayaran Sukses",
				Description: &desc,
			}
			metadata, _ := json.Marshal(map[string]any{
				"payment_id":      payment.ID.String(),
				"transaction_ref": payment.TransactionRef,
				"amount":          payment.Amount,
			})
			timeline.Metadata = metadata
			if err := tx.Create(&timeline).Error; err != nil {
				return err
			}

			var order entity.Order
			if err := tx.First(&order, "id = ?", payment.OrderID).Error; err == nil {
				var wallet entity.WorkerWallet
				errWallet := tx.First(&wallet, "worker_id = ?", order.WorkerID).Error
				if errWallet == nil {
					balanceBefore := wallet.Balance
					earningAmount := int64(payment.Amount)
					balanceAfter := balanceBefore + earningAmount

					if err := tx.Model(&wallet).Updates(map[string]any{
						"balance":        balanceAfter,
						"total_earnings": wallet.TotalEarnings + earningAmount,
					}).Error; err != nil {
						return err
					}

					txDesc := fmt.Sprintf("Pendapatan jasa dari Order %s", order.OrderNumber)
					walletTx := entity.WalletTransaction{
						WalletID:      wallet.ID,
						OrderID:       &order.ID,
						Type:          entity.WalletTxTypeEarning,
						Amount:        payment.Amount,
						BalanceBefore: balanceBefore,
						BalanceAfter:  balanceAfter,
						Description:   &txDesc,
						ReferenceID:   payment.TransactionRef,
						Status:        entity.WalletTxStatusCompleted,
						CompletedAt:   &now,
					}
					if err := tx.Create(&walletTx).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]any{
		"payment_id":     payment.ID.String(),
		"order_id":       payment.OrderID.String(),
		"payment_status": payment.PaymentStatus,
		"paid_at":        payment.PaidAt,
	}, nil
}

func (s *orderService) GetPaymentDetails(ctx context.Context, paymentID string) (any, error) {
	payID, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, http_error.BAD_REQUEST_ERROR
	}

	var payment entity.Payment
	if err := s.db.WithContext(ctx).First(&payment, "id = ?", payID).Error; err != nil {
		return nil, err
	}

	var invoice entity.Invoice
	_ = s.db.WithContext(ctx).First(&invoice, "id = ?", payment.InvoiceID).Error

	var order entity.Order
	_ = s.db.WithContext(ctx).First(&order, "id = ?", payment.OrderID).Error

	var customer entity.User
	_ = s.db.WithContext(ctx).First(&customer, "id = ?", payment.UserID).Error

	return map[string]any{
		"payment_id":     payment.ID.String(),
		"order_number":   order.OrderNumber,
		"order_title":    order.Title,
		"invoice_number": invoice.InvoiceNumber,
		"amount":         payment.Amount,
		"currency":       payment.Currency,
		"payment_status": payment.PaymentStatus,
		"payment_method": payment.PaymentMethod,
		"customer_name":  customer.FullName,
		"customer_email": customer.Email,
		"created_at":     payment.CreatedAt,
	}, nil
}

func (s *orderService) DownloadInvoicePDF(ctx context.Context, orderID string) ([]byte, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.getOrderForUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}
	content := fmt.Sprintf("Invoice %s\nOrder: %s\n", order.OrderNumber, order.ID.String())
	return []byte(content), nil
}

func (s *orderService) getOrderForUser(ctx context.Context, orderID string, userID uuid.UUID) (entity.Order, error) {
	id, err := parseUUID(orderID)
	if err != nil {
		return entity.Order{}, err
	}
	var order entity.Order
	err = s.db.WithContext(ctx).First(&order, "id = ? AND user_id = ?", id, userID).Error
	return order, err
}

func (s *orderService) getPurchase(ctx context.Context, orderID uuid.UUID, purchaseID string) (entity.Purchase, error) {
	id, err := parseUUID(purchaseID)
	if err != nil {
		return entity.Purchase{}, err
	}
	var purchase entity.Purchase
	err = s.db.WithContext(ctx).First(&purchase, "id = ? AND order_id = ?", id, orderID).Error
	return purchase, err
}

func (s *orderService) orderList(ctx context.Context, orders []entity.Order) []map[string]any {
	data := make([]map[string]any, 0, len(orders))
	for _, order := range orders {
		item, _ := s.orderDetail(ctx, order)
		data = append(data, item)
	}
	return data
}

func (s *orderService) orderDetail(ctx context.Context, order entity.Order) (map[string]any, error) {
	var worker entity.User
	var service entity.Service
	var category entity.Category
	_ = s.db.WithContext(ctx).First(&worker, "id = ?", order.WorkerID).Error
	_ = s.db.WithContext(ctx).First(&service, "id = ?", order.ServiceID).Error
	_ = s.db.WithContext(ctx).First(&category, "id = ?", order.CategoryID).Error
	return map[string]any{
		"order_id":              order.ID.String(),
		"order_number":          order.OrderNumber,
		"user_id":               order.UserID.String(),
		"worker_id":             order.WorkerID.String(),
		"worker_name":           worker.FullName,
		"service_id":            order.ServiceID.String(),
		"service_name":          service.Name,
		"category_id":           order.CategoryID.String(),
		"category_name":         category.Name,
		"title":                 order.Title,
		"description":           order.Description,
		"status":                order.Status,
		"urgency":               order.Urgency,
		"location_address":      order.LocationAddress,
		"location_detail":       order.LocationDetail,
		"location_lat":          order.LocationLat,
		"location_lng":          order.LocationLng,
		"preferred_date":        order.PreferredDate,
		"preferred_time_start":  order.PreferredTimeStart,
		"preferred_time_end":    order.PreferredTimeEnd,
		"notes":                 order.Notes,
		"booking_fee":           order.BookingFee,
		"base_service_fee":      order.BaseServiceFee,
		"total_material_cost":   order.TotalMaterialCost,
		"total_additional_cost": order.TotalAdditionalCost,
		"grand_total":           order.GrandTotal,
		"cancellation_reason":   order.CancellationReason,
		"cancellation_category": order.CancellationCategory,
		"accepted_at":           order.AcceptedAt,
		"started_at":            order.StartedAt,
		"completed_at":          order.CompletedAt,
		"cancelled_at":          order.CancelledAt,
		"created_at":            order.CreatedAt,
		"updated_at":            order.UpdatedAt,
	}, nil
}

func purchaseResponse(purchase entity.Purchase) map[string]any {
	return map[string]any{
		"purchase_id":            purchase.ID.String(),
		"order_id":               purchase.OrderID.String(),
		"worker_id":              purchase.WorkerID.String(),
		"item_name":              purchase.ItemName,
		"category":               purchase.Category,
		"quantity":               purchase.Quantity,
		"unit":                   purchase.Unit,
		"unit_price":             purchase.UnitPrice,
		"total_price":            purchase.TotalPrice,
		"reason":                 purchase.Reason,
		"receipt_photo_url":      purchase.ReceiptPhotoURL,
		"status":                 purchase.Status,
		"confidence":             purchase.Confidence,
		"needs_clarification":    purchase.NeedsClarification,
		"clarification_question": purchase.ClarificationQuestion,
		"clarification_response": purchase.ClarificationResponse,
		"created_at":             purchase.CreatedAt,
		"updated_at":             purchase.UpdatedAt,
	}
}

func chatResponse(message entity.ChatMessage) map[string]any {
	return map[string]any{
		"message_id":   message.ID.String(),
		"order_id":     message.OrderID.String(),
		"sender_id":    message.SenderID.String(),
		"sender_type":  message.SenderType,
		"content":      message.Content,
		"message_type": message.MessageType,
		"media_url":    message.MediaURL,
		"is_read":      message.IsRead,
		"read_at":      message.ReadAt,
		"created_at":   message.CreatedAt,
	}
}

func reviewResponse(review entity.Review) map[string]any {
	return map[string]any{
		"review_id":   review.ID.String(),
		"order_id":    review.OrderID.String(),
		"user_id":     review.UserID.String(),
		"worker_id":   review.WorkerID.String(),
		"review_type": review.ReviewType,
		"rating":      review.Rating,
		"comment":     review.Comment,
		"created_at":  review.CreatedAt,
	}
}

func (s *orderService) createReviewWithTags(ctx context.Context, review entity.Review, tags []string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing entity.Review
		if err := tx.First(&existing, "order_id = ? AND review_type = ?", review.OrderID, review.ReviewType).Error; err == nil {
			return http_error.DUPLICATE_DATA
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Create(&review).Error; err != nil {
			return err
		}
		for _, tag := range tags {
			if tag == "" {
				continue
			}
			if err := tx.Create(&entity.ReviewTag{ReviewID: review.ID, Tag: tag}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *orderService) invoiceResponse(ctx context.Context, orderID uuid.UUID) (map[string]any, error) {
	var invoice entity.Invoice
	if err := s.db.WithContext(ctx).First(&invoice, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	var lineItems []entity.InvoiceLineItem
	_ = s.db.WithContext(ctx).Where("invoice_id = ?", invoice.ID).Order("display_order ASC").Find(&lineItems).Error
	return map[string]any{
		"invoice_id":             invoice.ID.String(),
		"order_id":               invoice.OrderID.String(),
		"invoice_number":         invoice.InvoiceNumber,
		"base_service_fee":       invoice.BaseServiceFee,
		"total_material_cost":    invoice.TotalMaterialCost,
		"total_additional_cost":  invoice.TotalAdditionalCost,
		"booking_fee":            invoice.BookingFee,
		"platform_fee":           invoice.PlatformFee,
		"discount_amount":        invoice.DiscountAmount,
		"promo_code":             invoice.PromoCode,
		"grand_total":            invoice.GrandTotal,
		"currency":               invoice.Currency,
		"payment_instruction":    invoice.PaymentInstruction,
		"ai_work_summary":        invoice.AIWorkSummary,
		"ai_materials_summary":   invoice.AIMaterialsSummary,
		"worker_notes":           invoice.WorkerNotes,
		"all_purchases_approved": invoice.AllPurchasesApproved,
		"line_items":             lineItems,
		"created_at":             invoice.CreatedAt,
	}, nil
}

func createTimeline(tx *gorm.DB, orderID uuid.UUID, event string, label string, description *string, actorID *uuid.UUID, actorType string) error {
	timeline := entity.OrderTimeline{
		OrderID:     orderID,
		Event:       event,
		Label:       label,
		Description: description,
		ActorID:     actorID,
		ActorType:   &actorType,
	}
	metadata, _ := json.Marshal(map[string]any{})
	timeline.Metadata = metadata
	return tx.Create(&timeline).Error
}

func parseOptionalDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, http_error.BAD_REQUEST_ERROR
	}
	return &parsed, nil
}

func parseOptionalClock(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("15:04", *value)
	if err != nil {
		return nil, http_error.BAD_REQUEST_ERROR
	}
	return &parsed, nil
}

func parseUUIDList(values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		id, err := parseUUID(value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func newOrderNumber() string {
	return fmt.Sprintf("HD-%s-%s", time.Now().Format("20060102"), uuid.NewString()[:8])
}

func estimateETA(distance float64) int {
	if distance <= 0 {
		return 0
	}
	return int(mathCeil(distance / 25 * 60))
}

func mathCeil(value float64) float64 {
	if value == float64(int(value)) {
		return value
	}
	return float64(int(value) + 1)
}
