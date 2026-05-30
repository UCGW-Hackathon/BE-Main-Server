package entity

import "encoding/json"

type JSONB = json.RawMessage

type UserRole string

const (
	UserRoleUser   UserRole = "user"
	UserRoleWorker UserRole = "worker"
	UserRoleAdmin  UserRole = "admin"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusOnTheWay   OrderStatus = "on_the_way"
	OrderStatusArrived    OrderStatus = "arrived"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusWorkPaused OrderStatus = "work_paused"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRejected   OrderStatus = "rejected"
)

type OrderUrgency string

const (
	OrderUrgencyNormal OrderUrgency = "normal"
	OrderUrgencyUrgent OrderUrgency = "urgent"
)

type PurchaseStatus string

const (
	PurchaseStatusDraft              PurchaseStatus = "draft"
	PurchaseStatusPendingApproval    PurchaseStatus = "pending_approval"
	PurchaseStatusApproved           PurchaseStatus = "approved"
	PurchaseStatusRejected           PurchaseStatus = "rejected"
	PurchaseStatusNeedsClarification PurchaseStatus = "needs_clarification"
)

type PurchaseCategory string

const (
	PurchaseCategoryMaterial      PurchaseCategory = "material"
	PurchaseCategoryAlat          PurchaseCategory = "alat"
	PurchaseCategorySparepart     PurchaseCategory = "sparepart"
	PurchaseCategoryBahanBangunan PurchaseCategory = "bahan_bangunan"
	PurchaseCategoryBiayaTambahan PurchaseCategory = "biaya_tambahan"
	PurchaseCategoryLainnya       PurchaseCategory = "lainnya"
)

type RiskFlagType string

const (
	RiskFlagTypeHargaTidakWajar    RiskFlagType = "harga_tidak_wajar"
	RiskFlagTypeItemTidakRelevan   RiskFlagType = "item_tidak_relevan"
	RiskFlagTypeDataTidakLengkap   RiskFlagType = "data_tidak_lengkap"
	RiskFlagTypeNotaTidakJelas     RiskFlagType = "nota_tidak_jelas"
	RiskFlagTypeDuplikat           RiskFlagType = "duplikat"
	RiskFlagTypeAlasanTidakLengkap RiskFlagType = "alasan_tidak_lengkap"
)

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeSystem MessageType = "system"
)

type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "ewallet"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid   PaymentStatus = "unpaid"
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type VerificationStatus string

const (
	VerificationStatusUnverified VerificationStatus = "unverified"
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusRejected   VerificationStatus = "rejected"
)

type NotificationType string

const (
	NotificationTypeOrder    NotificationType = "order"
	NotificationTypePurchase NotificationType = "purchase"
	NotificationTypeChat     NotificationType = "chat"
	NotificationTypePromo    NotificationType = "promo"
	NotificationTypeSystem   NotificationType = "system"
	NotificationTypePayment  NotificationType = "payment"
)

type ArticleCategory string

const (
	ArticleCategoryFAQ     ArticleCategory = "faq"
	ArticleCategoryGuide   ArticleCategory = "guide"
	ArticleCategoryTips    ArticleCategory = "tips"
	ArticleCategorySafety  ArticleCategory = "safety"
	ArticleCategoryPayment ArticleCategory = "payment"
)

type FAQCategory string

const (
	FAQCategoryGeneral      FAQCategory = "general"
	FAQCategoryPayment      FAQCategory = "payment"
	FAQCategoryTracking     FAQCategory = "tracking"
	FAQCategorySecurity     FAQCategory = "security"
	FAQCategoryCancellation FAQCategory = "cancellation"
)

type WalletTxType string

const (
	WalletTxTypeEarning    WalletTxType = "earning"
	WalletTxTypeWithdrawal WalletTxType = "withdrawal"
	WalletTxTypeRefund     WalletTxType = "refund"
	WalletTxTypeBonus      WalletTxType = "bonus"
	WalletTxTypeFee        WalletTxType = "fee"
)

type WalletTxStatus string

const (
	WalletTxStatusPending   WalletTxStatus = "pending"
	WalletTxStatusCompleted WalletTxStatus = "completed"
	WalletTxStatusFailed    WalletTxStatus = "failed"
	WalletTxStatusCancelled WalletTxStatus = "cancelled"
)

type AuditAction string

const (
	AuditActionCreated                AuditAction = "created"
	AuditActionAIProcessed            AuditAction = "ai_processed"
	AuditActionSubmitted              AuditAction = "submitted"
	AuditActionApproved               AuditAction = "approved"
	AuditActionRejected               AuditAction = "rejected"
	AuditActionClarificationRequested AuditAction = "clarification_requested"
	AuditActionClarificationResponded AuditAction = "clarification_responded"
	AuditActionEdited                 AuditAction = "edited"
	AuditActionDeleted                AuditAction = "deleted"
)
