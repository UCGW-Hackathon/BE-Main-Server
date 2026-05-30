package controllers

import (
	"net/http"

	"whatsapp-backend/dto"
	"whatsapp-backend/services"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
)

type WorkerController interface {
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	UpdateCoverPhoto(ctx *gin.Context)
	SubmitVerification(ctx *gin.Context)
	GetVerification(ctx *gin.Context)
	GetHome(ctx *gin.Context)
	UpdateAvailability(ctx *gin.Context)
	ListIncomingOrders(ctx *gin.Context)
	GetIncomingOrderDetail(ctx *gin.Context)
	AcceptOrder(ctx *gin.Context)
	RejectOrder(ctx *gin.Context)
	ListOrders(ctx *gin.Context)
	GetOrderDetail(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	GenerateInvoice(ctx *gin.Context)
	AddPurchase(ctx *gin.Context)
	AIProcessPurchase(ctx *gin.Context)
	ReceiptScanPurchase(ctx *gin.Context)
	UpdatePurchase(ctx *gin.Context)
	DeletePurchase(ctx *gin.Context)
	SubmitPurchase(ctx *gin.Context)
	BulkSubmitPurchase(ctx *gin.Context)
	ClarifyPurchaseResponse(ctx *gin.Context)
	ListChatMessages(ctx *gin.Context)
	SendChatMessage(ctx *gin.Context)
	MarkChatRead(ctx *gin.Context)
	ListChats(ctx *gin.Context)
	CreateCustomerRating(ctx *gin.Context)
	GetCustomerRating(ctx *gin.Context)
	GetHistory(ctx *gin.Context)
	GetStatistics(ctx *gin.Context)
	GetWallet(ctx *gin.Context)
	ListWalletTransactions(ctx *gin.Context)
	Withdraw(ctx *gin.Context)
	UpdateLocation(ctx *gin.Context)
}

type workerController struct {
	workerService services.WorkerService
}

func NewWorkerController(workerService services.WorkerService) WorkerController {
	return &workerController{workerService: workerService}
}

func (wc *workerController) GetProfile(ctx *gin.Context) {
	data, err := wc.workerService.GetProfile(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) UpdateProfile(ctx *gin.Context) {
	req := RequestJSON[dto.WorkerProfileUpdateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.UpdateProfile(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Profil worker berhasil diperbarui", data, nil)
}

func (wc *workerController) UpdateCoverPhoto(ctx *gin.Context) {
	file, err := ctx.FormFile("cover_photo")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Cover photo is required", nil)
		return
	}

	coverURL, err := utils.SaveUploadedFile(ctx, file, "worker-cover")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	data, err := wc.workerService.UpdateCoverPhoto(ctx, coverURL)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) SubmitVerification(ctx *gin.Context) {
	idCard, err := ctx.FormFile("id_card")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "id_card is required", nil)
		return
	}

	idCardURL, err := utils.SaveUploadedFile(ctx, idCard, "verification")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	data, err := wc.workerService.SubmitVerification(ctx, idCardURL, nil)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Dokumen verifikasi berhasil dikirim", data, nil)
}

func (wc *workerController) GetVerification(ctx *gin.Context) {
	data, err := wc.workerService.GetVerification(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) GetHome(ctx *gin.Context) {
	data, err := wc.workerService.GetHome(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) UpdateAvailability(ctx *gin.Context) {
	req := RequestJSON[dto.WorkerAvailabilityRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.UpdateAvailability(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Status ketersediaan berhasil diperbarui", data, nil)
}

func (wc *workerController) ListIncomingOrders(ctx *gin.Context) {
	data, err := wc.workerService.ListIncomingOrders(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) GetIncomingOrderDetail(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := wc.workerService.GetIncomingOrderDetail(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) AcceptOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerOrderAcceptRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.AcceptOrder(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Order berhasil diterima", data, nil)
}

func (wc *workerController) RejectOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerOrderRejectRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.RejectOrder(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Order ditolak", data, nil)
}

func (wc *workerController) ListOrders(ctx *gin.Context) {
	data, err := wc.workerService.ListOrders(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     10,
		"total":        0,
		"total_pages":  0,
	})
}

func (wc *workerController) GetOrderDetail(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := wc.workerService.GetOrderDetail(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) UpdateOrderStatus(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerOrderStatusRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.UpdateOrderStatus(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Status order berhasil diperbarui", data, nil)
}

func (wc *workerController) GenerateInvoice(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerGenerateInvoiceRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.GenerateInvoice(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Invoice berhasil dibuat", data, nil)
}

func (wc *workerController) AddPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerPurchaseCreateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.AddPurchase(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Pembelian berhasil ditambahkan", data, nil)
}

func (wc *workerController) AIProcessPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerPurchaseAIProcessRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.AIProcessPurchase(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Input berhasil diproses oleh AI", data, nil)
}

func (wc *workerController) ReceiptScanPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	file, err := ctx.FormFile("receipt_photo")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "receipt_photo is required", nil)
		return
	}

	receiptURL, err := utils.SaveUploadedFile(ctx, file, "receipts")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	data, err := wc.workerService.ReceiptScanPurchase(ctx, orderID, receiptURL)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Nota berhasil dipindai", data, nil)
}

func (wc *workerController) UpdatePurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	req := RequestJSON[dto.WorkerPurchaseUpdateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.UpdatePurchase(ctx, orderID, purchaseID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil diperbarui", data, nil)
}

func (wc *workerController) DeletePurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	err := wc.workerService.DeletePurchase(ctx, orderID, purchaseID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil dihapus", gin.H{}, nil)
}

func (wc *workerController) SubmitPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	data, err := wc.workerService.SubmitPurchase(ctx, orderID, purchaseID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil dikirim", data, nil)
}

func (wc *workerController) BulkSubmitPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.WorkerPurchaseBulkSubmitRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.BulkSubmitPurchase(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil dikirim", data, nil)
}

func (wc *workerController) ClarifyPurchaseResponse(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	req := RequestJSON[dto.WorkerPurchaseClarifyResponseRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.ClarifyPurchaseResponse(ctx, orderID, purchaseID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Klarifikasi berhasil dikirim", data, nil)
}

func (wc *workerController) ListChatMessages(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := wc.workerService.ListChatMessages(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) SendChatMessage(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.ChatSendRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.SendChatMessage(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "", data, nil)
}

func (wc *workerController) MarkChatRead(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	err := wc.workerService.MarkChatRead(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Semua pesan ditandai sudah dibaca", gin.H{}, nil)
}

func (wc *workerController) ListChats(ctx *gin.Context) {
	data, err := wc.workerService.ListChats(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     20,
		"total":        0,
		"total_pages":  0,
	})
}

func (wc *workerController) CreateCustomerRating(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.CustomerRatingCreateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.CreateCustomerRating(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Rating konsumen berhasil dikirim", data, nil)
}

func (wc *workerController) GetCustomerRating(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := wc.workerService.GetCustomerRating(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) GetHistory(ctx *gin.Context) {
	data, err := wc.workerService.GetHistory(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     10,
		"total":        0,
		"total_pages":  0,
	})
}

func (wc *workerController) GetStatistics(ctx *gin.Context) {
	data, err := wc.workerService.GetStatistics(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) GetWallet(ctx *gin.Context) {
	data, err := wc.workerService.GetWallet(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerController) ListWalletTransactions(ctx *gin.Context) {
	data, err := wc.workerService.ListWalletTransactions(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     20,
		"total":        0,
		"total_pages":  0,
	})
}

func (wc *workerController) Withdraw(ctx *gin.Context) {
	req := RequestJSON[dto.WalletWithdrawRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.Withdraw(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Permintaan penarikan berhasil dibuat", data, nil)
}

func (wc *workerController) UpdateLocation(ctx *gin.Context) {
	req := RequestJSON[dto.WorkerLocationUpdateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := wc.workerService.UpdateLocation(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}
