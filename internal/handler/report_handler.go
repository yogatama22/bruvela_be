package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bruvela-backend/internal/repository"
)

type ReportHandler struct {
	orderRepo      repository.OrderRepository
	financeRepo    repository.FinanceRepository
	ingredientRepo repository.IngredientRepository
	recipeRepo     repository.RecipeRepository
	productRepo    repository.ProductRepository
}

func NewReportHandler(
	orderRepo repository.OrderRepository,
	financeRepo repository.FinanceRepository,
	ingredientRepo repository.IngredientRepository,
	recipeRepo repository.RecipeRepository,
	productRepo repository.ProductRepository,
) *ReportHandler {
	return &ReportHandler{
		orderRepo:      orderRepo,
		financeRepo:    financeRepo,
		ingredientRepo: ingredientRepo,
		recipeRepo:     recipeRepo,
		productRepo:    productRepo,
	}
}

// setCSVHeaders writes standard CSV download headers.
func setCSVHeaders(c *gin.Context, filename string) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
}

// =========================
// 1. Orders Export
// =========================
func (h *ReportHandler) ExportOrders(c *gin.Context) {
	batchIDStr := c.Query("batch_id")
	if batchIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id required"})
		return
	}
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch_id"})
		return
	}

	orders, _, err := h.orderRepo.FindAll(map[string]interface{}{"batch_id": batchID}, 10000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	setCSVHeaders(c, fmt.Sprintf("orders_batch_%s_%s.csv", batchIDStr[:8], time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Header row
	writer.Write([]string{
		"Order ID", "Tanggal", "Customer", "Channel", "Shipping Type",
		"Shipping Dest", "Shipping Cost", "Discount", "Total Product",
		"Total Bill", "Pay Status", "Prod Status", "Note", "Items",
	})

	for _, order := range orders {
		itemsStr := ""
		for i, item := range order.Items {
			if i > 0 {
				itemsStr += "; "
			}
			name := item.ProductName
			if name == "" && item.Product != nil {
				name = item.Product.Name
			}
			itemsStr += fmt.Sprintf("%s x%d", name, item.QtyBox)
		}
		writer.Write([]string{
			order.ID.String(),
			order.OrderDate.Format("2006-01-02"),
			order.CustomerName,
			order.Channel,
			order.ShippingType,
			order.ShippingDest,
			strconv.Itoa(order.ShippingCost),
			strconv.Itoa(order.Discount),
			strconv.Itoa(order.TotalProduct),
			strconv.Itoa(order.TotalBill),
			order.PayStatus,
			order.ProdStatus,
			order.Note,
			itemsStr,
		})
	}
}

// =========================
// 2. Finance/Journal Export
// =========================
func (h *ReportHandler) ExportFinance(c *gin.Context) {
	batchIDStr := c.Query("batch_id")
	if batchIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "batch_id required"})
		return
	}
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch_id"})
		return
	}

	entries, err := h.financeRepo.FindJournalEntries(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	setCSVHeaders(c, fmt.Sprintf("journal_batch_%s_%s.csv", batchIDStr[:8], time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{
		"Tanggal", "Tipe", "Keterangan", "Partner", "Jumlah", "Saldo",
	})

	runningBalance := 0
	for _, entry := range entries {
		runningBalance += entry.Amount
		writer.Write([]string{
			entry.EntryDate.Format("2006-01-02"),
			entry.Type,
			entry.Description,
			entry.Partner,
			strconv.Itoa(entry.Amount),
			strconv.Itoa(runningBalance),
		})
	}
}

// =========================
// 3. Inventory Export
// =========================
func (h *ReportHandler) ExportInventory(c *gin.Context) {
	ingredients, err := h.ingredientRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	setCSVHeaders(c, fmt.Sprintf("inventory_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{
		"Nama Bahan", "Satuan Pack", "Qty/Pack", "Satuan Pakai",
		"Harga/Pack", "Harga/Satuan Pakai", "Stok Minimum",
		"Stok Saat Ini", "Status",
	})

	for _, ing := range ingredients {
		status := "AMAN"
		if ing.CurrentStock < 0 {
			status = "MINUS"
		} else if ing.CurrentStock < ing.MinStock {
			status = "KRITIS"
		}
		writer.Write([]string{
			ing.Name,
			ing.PackUnit,
			fmt.Sprintf("%.3f", ing.QtyPerPack),
			ing.UseUnit,
			strconv.Itoa(ing.PricePerPack),
			fmt.Sprintf("%.4f", ing.PricePerUse),
			fmt.Sprintf("%.3f", ing.MinStock),
			fmt.Sprintf("%.3f", ing.CurrentStock),
			status,
		})
	}
}

// =========================
// 4. HPP Report (Recipe + Cost)
// =========================
func (h *ReportHandler) ExportHPP(c *gin.Context) {
	products, err := h.productRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	setCSVHeaders(c, fmt.Sprintf("hpp_report_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{
		"Kode", "Nama Produk", "Harga Jual", "Pcs/Box", "HPP/Box",
		"Laba/Box", "Margin %", "Status",
	})

	for _, product := range products {
		recipes, _ := h.recipeRepo.FindByProductID(product.ID)
		hpp := 0.0
		for _, r := range recipes {
			hpp += r.CostPerBox
		}
		profit := float64(product.Price) - hpp
		margin := 0.0
		if product.Price > 0 {
			margin = (profit / float64(product.Price)) * 100
		}
		writer.Write([]string{
			product.Code,
			product.Name,
			strconv.Itoa(product.Price),
			strconv.Itoa(product.PcsPerBox),
			fmt.Sprintf("%.2f", hpp),
			fmt.Sprintf("%.2f", profit),
			fmt.Sprintf("%.2f", margin),
			product.Status,
		})
	}
}
