package handler

import (
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	orderRepo      repository.OrderRepository
	ingredientRepo repository.IngredientRepository
	batchRepo      repository.BatchRepository
}

func NewDashboardHandler(
	orderRepo repository.OrderRepository,
	ingredientRepo repository.IngredientRepository,
	batchRepo repository.BatchRepository,
) *DashboardHandler {
	return &DashboardHandler{
		orderRepo:      orderRepo,
		ingredientRepo: ingredientRepo,
		batchRepo:      batchRepo,
	}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
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

	filters := map[string]interface{}{
		"batch_id": batchID,
	}
	orders, total, err := h.orderRepo.FindAll(filters, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalRevenue := 0
	totalPaid := 0
	totalPending := 0
	statusCount := map[string]int{
		"baru":       0,
		"proses":     0,
		"siap_kirim": 0,
		"selesai":    0,
		"batal":      0,
	}

	for _, order := range orders {
		totalRevenue += order.TotalBill
		if order.PayStatus == "lunas" {
			totalPaid += order.TotalBill
		} else {
			totalPending += order.TotalBill
		}
		statusCount[order.ProdStatus]++
	}

	lowStockIngredients, _ := h.ingredientRepo.FindLowStock()

	c.JSON(http.StatusOK, gin.H{
		"total_orders":          total,
		"total_revenue":         totalRevenue,
		"total_paid":            totalPaid,
		"total_pending":         totalPending,
		"status_count":          statusCount,
		"low_stock_count":       len(lowStockIngredients),
		"low_stock_ingredients": lowStockIngredients,
	})
}
