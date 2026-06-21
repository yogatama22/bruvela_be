package handler

import (
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockLogHandler struct {
	stockLogRepo repository.StockLogRepository
}

func NewStockLogHandler(stockLogRepo repository.StockLogRepository) *StockLogHandler {
	return &StockLogHandler{stockLogRepo: stockLogRepo}
}

func (h *StockLogHandler) GetAll(c *gin.Context) {
	filters := make(map[string]interface{})

	if batchID := c.Query("batch_id"); batchID != "" {
		filters["batch_id"] = batchID
	}
	if ingredientID := c.Query("ingredient_id"); ingredientID != "" {
		filters["ingredient_id"] = ingredientID
	}
	if logType := c.Query("log_type"); logType != "" {
		filters["log_type"] = logType
	}
	if referenceType := c.Query("reference_type"); referenceType != "" {
		filters["reference_type"] = referenceType
	}

	logs, err := h.stockLogRepo.FindAll(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
