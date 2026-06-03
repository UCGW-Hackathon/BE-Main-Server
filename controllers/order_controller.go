package controllers

import (
	"fmt"
	"net/http"

	"situkang/dto"
	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	ListOrders(ctx *gin.Context)
	GetOrderDetail(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	GetTracking(ctx *gin.Context)
	GetTrackingLocation(ctx *gin.Context)
	ListPurchases(ctx *gin.Context)
	GetPurchaseDetail(ctx *gin.Context)
	ApprovePurchase(ctx *gin.Context)
	RejectPurchase(ctx *gin.Context)
	ClarifyPurchase(ctx *gin.Context)
	BulkApprovePurchases(ctx *gin.Context)
	ListChatMessages(ctx *gin.Context)
	SendChatMessage(ctx *gin.Context)
	MarkChatRead(ctx *gin.Context)
	ListChats(ctx *gin.Context)
	CreateRating(ctx *gin.Context)
	GetRating(ctx *gin.Context)
	GetInvoice(ctx *gin.Context)
	CreatePayment(ctx *gin.Context)
	DownloadInvoicePDF(ctx *gin.Context)
	SandboxCheckout(ctx *gin.Context)
	SandboxCallback(ctx *gin.Context)
}

type orderController struct {
	orderService services.OrderService
}

func NewOrderController(orderService services.OrderService) OrderController {
	return &orderController{orderService: orderService}
}

func (oc *orderController) CreateOrder(ctx *gin.Context) {
	req := RequestJSON[dto.OrderCreateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.CreateOrder(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Pesanan berhasil dibuat", data, nil)
}

func (oc *orderController) ListOrders(ctx *gin.Context) {
	data, err := oc.orderService.ListOrders(ctx)
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

func (oc *orderController) GetOrderDetail(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.GetOrderDetail(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) CancelOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.OrderCancelRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.CancelOrder(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pesanan berhasil dibatalkan", data, nil)
}

func (oc *orderController) GetTracking(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.GetTracking(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) GetTrackingLocation(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.GetTrackingLocation(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) ListPurchases(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.ListPurchases(ctx, orderID)
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

func (oc *orderController) GetPurchaseDetail(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	data, err := oc.orderService.GetPurchaseDetail(ctx, orderID, purchaseID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) ApprovePurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	req := RequestJSON[dto.PurchaseApproveRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.ApprovePurchase(ctx, orderID, purchaseID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil disetujui", data, nil)
}

func (oc *orderController) RejectPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	req := RequestJSON[dto.PurchaseRejectRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.RejectPurchase(ctx, orderID, purchaseID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian ditolak", data, nil)
}

func (oc *orderController) ClarifyPurchase(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	purchaseID := ctx.Param("purchase_id")
	req := RequestJSON[dto.PurchaseClarifyRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.ClarifyPurchase(ctx, orderID, purchaseID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Permintaan klarifikasi berhasil dikirim", data, nil)
}

func (oc *orderController) BulkApprovePurchases(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.PurchaseBulkApproveRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.BulkApprovePurchases(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembelian berhasil disetujui", data, nil)
}

func (oc *orderController) ListChatMessages(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.ListChatMessages(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) SendChatMessage(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.ChatSendRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.SendChatMessage(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "", data, nil)
}

func (oc *orderController) MarkChatRead(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	err := oc.orderService.MarkChatRead(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Semua pesan ditandai sudah dibaca", gin.H{}, nil)
}

func (oc *orderController) ListChats(ctx *gin.Context) {
	data, err := oc.orderService.ListChats(ctx)
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

func (oc *orderController) CreateRating(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.RatingCreateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.CreateRating(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Rating berhasil dikirim", data, nil)
}

func (oc *orderController) GetRating(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.GetRating(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) GetInvoice(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	data, err := oc.orderService.GetInvoice(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (oc *orderController) CreatePayment(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	req := RequestJSON[dto.PaymentCreateRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := oc.orderService.CreatePayment(ctx, orderID, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Pembayaran berhasil dikonfirmasi", data, nil)
}

func (oc *orderController) DownloadInvoicePDF(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	content, err := oc.orderService.DownloadInvoicePDF(ctx, orderID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=\"invoice.pdf\"")
	ctx.Data(http.StatusOK, "application/pdf", content)
}

func (oc *orderController) SandboxCheckout(ctx *gin.Context) {
	paymentID := ctx.Query("payment_id")
	if paymentID == "" {
		ctx.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte("<h1>Bad Request</h1><p>Missing payment_id parameter</p>"))
		return
	}

	details, err := oc.orderService.GetPaymentDetails(ctx, paymentID)
	if err != nil {
		ctx.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(fmt.Sprintf("<h1>Not Found</h1><p>%s</p>", err.Error())))
		return
	}

	data := details.(map[string]any)
	amountFormatted := fmt.Sprintf("Rp %s", formatRupiah(data["amount"].(int)))

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SiTukang - Sandbox Checkout</title>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
            font-family: 'Plus Jakarta Sans', sans-serif;
        }
        body {
            background: linear-gradient(135deg, #0f172a 0%%, #1e1b4b 100%%);
            color: #f1f5f9;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: rgba(30, 41, 59, 0.7);
            backdrop-filter: blur(16px);
            border: 1px solid rgba(255, 255, 255, 0.08);
            border-radius: 24px;
            width: 100%%;
            max-width: 500px;
            padding: 40px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
            text-align: center;
        }
        .logo {
            font-size: 28px;
            font-weight: 700;
            background: linear-gradient(135deg, #38bdf8 0%%, #818cf8 100%%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 8px;
        }
        .subtitle {
            font-size: 14px;
            color: #94a3b8;
            margin-bottom: 32px;
            text-transform: uppercase;
            letter-spacing: 2px;
            font-weight: 600;
        }
        .amount-card {
            background: rgba(15, 23, 42, 0.6);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 32px;
        }
        .amount-label {
            font-size: 12px;
            color: #94a3b8;
            margin-bottom: 8px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .amount-value {
            font-size: 32px;
            font-weight: 700;
            color: #38bdf8;
        }
        .details-list {
            text-align: left;
            margin-bottom: 36px;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            padding: 12px 0;
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        .detail-row:last-child {
            border-bottom: none;
        }
        .detail-label {
            color: #94a3b8;
            font-size: 14px;
        }
        .detail-value {
            font-weight: 500;
            font-size: 14px;
            color: #f1f5f9;
        }
        .pay-btn {
            background: linear-gradient(135deg, #0284c7 0%%, #4f46e5 100%%);
            color: white;
            border: none;
            border-radius: 12px;
            padding: 16px 32px;
            font-size: 16px;
            font-weight: 600;
            width: 100%%;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
        }
        .pay-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(79, 70, 229, 0.4);
            background: linear-gradient(135deg, #0ea5e9 0%%, #6366f1 100%%);
        }
        .pay-btn:active {
            transform: translateY(0);
        }
        .status-screen {
            display: none;
            animation: fadeIn 0.5s ease forwards;
        }
        .success-icon {
            width: 72px;
            height: 72px;
            background: rgba(16, 185, 129, 0.1);
            border: 2px solid #10b981;
            color: #10b981;
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 36px;
            margin: 0 auto 24px auto;
        }
        .success-title {
            font-size: 24px;
            font-weight: 700;
            margin-bottom: 12px;
            color: #10b981;
        }
        .success-desc {
            color: #94a3b8;
            font-size: 14px;
            line-height: 1.6;
            margin-bottom: 32px;
        }
        .close-btn {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            color: #f1f5f9;
            border-radius: 12px;
            padding: 12px 24px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .close-btn:hover {
            background: rgba(255, 255, 255, 0.1);
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: scale(0.95); }
            to { opacity: 1; transform: scale(1); }
        }
    </style>
</head>
<body>
    <div class="container" id="checkout-box">
        <div class="logo">SiTukang</div>
        <div class="subtitle">Sandbox Checkout</div>
        
        <div class="amount-card">
            <div class="amount-label">Total Tagihan</div>
            <div class="amount-value">%s</div>
        </div>

        <div class="details-list">
            <div class="detail-row">
                <span class="detail-label">Nomor Order</span>
                <span class="detail-value">%s</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Nama Pekerjaan</span>
                <span class="detail-value">%s</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Nomor Invoice</span>
                <span class="detail-value">%s</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Pelanggan</span>
                <span class="detail-value">%s</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Metode Pembayaran</span>
                <span class="detail-value" style="text-transform: uppercase;">%s</span>
            </div>
        </div>

        <button class="pay-btn" id="pay-button" onclick="simulatePayment()">Bayar Sekarang (Simulasi)</button>
    </div>

    <div class="container status-screen" id="success-box">
        <div class="success-icon">✓</div>
        <div class="success-title">Pembayaran Sukses!</div>
        <div class="success-desc">
            Transaksi sandbox berhasil diselesaikan. Status pesanan Anda telah otomatis diperbarui menjadi selesai (completed) dan pendapatan telah dikreditkan ke dompet worker.
        </div>
        <button class="close-btn" onclick="window.close()">Tutup Halaman</button>
    </div>

    <script>
        async function simulatePayment() {
            const button = document.getElementById('pay-button');
            button.disabled = true;
            button.innerText = 'Memproses...';
            button.style.opacity = '0.7';

            try {
                const response = await fetch('/v1/payments/sandbox-callback', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        payment_id: '%s',
                        status: 'success'
                    })
                });

                if (response.ok) {
                    document.getElementById('checkout-box').style.display = 'none';
                    document.getElementById('success-box').style.display = 'block';
                } else {
                    alert('Gagal mensimulasikan pembayaran. Silakan coba lagi.');
                    button.disabled = false;
                    button.innerText = 'Bayar Sekarang (Simulasi)';
                    button.style.opacity = '1';
                }
            } catch (err) {
                console.error(err);
                alert('Terjadi kesalahan jaringan.');
                button.disabled = false;
                button.innerText = 'Bayar Sekarang (Simulasi)';
                button.style.opacity = '1';
            }
        }
    </script>
</body>
</html>`, amountFormatted, data["order_number"], data["order_title"], data["invoice_number"], data["customer_name"], data["payment_method"], paymentID)

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (oc *orderController) SandboxCallback(ctx *gin.Context) {
	var req dto.SandboxPaymentCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}

	data, err := oc.orderService.SandboxCallback(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}

	utils.JSONSuccess(ctx, http.StatusOK, "Pembayaran sandbox sukses diproses", data, nil)
}

func formatRupiah(amount int) string {
	str := fmt.Sprintf("%d", amount)
	var result []rune
	count := 0
	for i := len(str) - 1; i >= 0; i-- {
		result = append([]rune{rune(str[i])}, result...)
		count++
		if count%3 == 0 && i != 0 {
			result = append([]rune{'.'}, result...)
		}
	}
	return string(result)
}
