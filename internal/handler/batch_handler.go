package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BatchHandler struct {
	batchRepo   repository.BatchRepository
	orderRepo   repository.OrderRepository
	financeRepo repository.FinanceRepository
	recipeRepo  repository.RecipeRepository
}

func NewBatchHandler(
	batchRepo repository.BatchRepository,
	orderRepo repository.OrderRepository,
	financeRepo repository.FinanceRepository,
	recipeRepo repository.RecipeRepository,
) *BatchHandler {
	return &BatchHandler{
		batchRepo:   batchRepo,
		orderRepo:   orderRepo,
		financeRepo: financeRepo,
		recipeRepo:  recipeRepo,
	}
}

func (h *BatchHandler) GetAll(c *gin.Context) {
	batches, err := h.batchRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, batches)
}

func (h *BatchHandler) GetActive(c *gin.Context) {
	batch, err := h.batchRepo.FindActive()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active batch found"})
		return
	}
	c.JSON(http.StatusOK, batch)
}

func (h *BatchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	batch, err := h.batchRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Batch not found"})
		return
	}

	c.JSON(http.StatusOK, batch)
}

func (h *BatchHandler) Create(c *gin.Context) {
	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if batch number already exists
	existing, _ := h.batchRepo.FindByBatchNumber(batch.BatchNumber)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Batch number already exists"})
		return
	}

	// If status is open, close all other batches
	if batch.Status == "open" {
		if err := h.batchRepo.CloseAllBatches(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close other batches"})
			return
		}
	}

	if err := h.batchRepo.Create(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, batch)
}

func (h *BatchHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch.ID = id

	// Check if batch number already exists (excluding current)
	existing, _ := h.batchRepo.FindByBatchNumber(batch.BatchNumber)
	if existing != nil && existing.ID != id {
		c.JSON(http.StatusConflict, gin.H{"error": "Batch number already exists"})
		return
	}

	// If status is open, close all other batches
	if batch.Status == "open" {
		if err := h.batchRepo.CloseAllBatches(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close other batches"})
			return
		}
	}

	if err := h.batchRepo.Update(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, batch)
}

func (h *BatchHandler) SetActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Close all batches first
	if err := h.batchRepo.CloseAllBatches(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close other batches"})
		return
	}

	// Open the selected batch
	if err := h.batchRepo.SetStatus(id, "open"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	batch, _ := h.batchRepo.FindByID(id)
	c.JSON(http.StatusOK, batch)
}

func (h *BatchHandler) Close(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.batchRepo.SetStatus(id, "closed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	batch, _ := h.batchRepo.FindByID(id)
	c.JSON(http.StatusOK, batch)
}

func (h *BatchHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.batchRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch deleted successfully"})
}

func (h *BatchHandler) GetSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	batch, err := h.batchRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Batch not found"})
		return
	}

	orders, err := h.orderRepo.FindByBatchID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ordersTotal := len(orders)
	revenue := 0
	totalPaid := 0
	hppTotal := 0.0

	for _, order := range orders {
		revenue += order.TotalBill
		if order.PayStatus == "lunas" {
			totalPaid += order.TotalBill
		}
		for _, item := range order.Items {
			recipes, _ := h.recipeRepo.FindByProductID(item.ProductID)
			for _, recipe := range recipes {
				hppTotal += recipe.CostPerBox * float64(item.QtyBox)
			}
		}
	}

	grossProfit := float64(revenue) - hppTotal
	marginPct := 0.0
	if revenue > 0 {
		marginPct = (grossProfit / float64(revenue)) * 100
	}

	financeSummary, _ := h.financeRepo.GetBatchSummary(id)

	c.JSON(http.StatusOK, gin.H{
		"batch":        batch,
		"orders_total": ordersTotal,
		"revenue":      revenue,
		"total_paid":   totalPaid,
		"hpp_total":    hppTotal,
		"gross_profit": grossProfit,
		"margin_pct":   marginPct,
		"finance":      financeSummary,
		"orders":       orders,
	})
}
