package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	orderRepo      repository.OrderRepository
	recipeRepo     repository.RecipeRepository
	ingredientRepo repository.IngredientRepository
	stockLogRepo   repository.StockLogRepository
	db             *gorm.DB
}

func NewOrderHandler(
	orderRepo repository.OrderRepository,
	recipeRepo repository.RecipeRepository,
	ingredientRepo repository.IngredientRepository,
	stockLogRepo repository.StockLogRepository,
	db *gorm.DB,
) *OrderHandler {
	return &OrderHandler{
		orderRepo:      orderRepo,
		recipeRepo:     recipeRepo,
		ingredientRepo: ingredientRepo,
		stockLogRepo:   stockLogRepo,
		db:             db,
	}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := make(map[string]interface{})
	if batchID := c.Query("batch_id"); batchID != "" {
		filters["batch_id"] = batchID
	}
	if payStatus := c.Query("pay_status"); payStatus != "" {
		filters["pay_status"] = payStatus
	}
	if prodStatus := c.Query("prod_status"); prodStatus != "" {
		filters["prod_status"] = prodStatus
	}

	orders, total, err := h.orderRepo.FindAll(filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   orders,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Safe cast to UUID
	if uid, ok := userID.(uuid.UUID); ok {
		order.CreatedBy = uid
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create order (totals will be calculated in repository after items are created)
	if err := h.orderRepo.Create(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Load existing order to preserve immutable fields
	existingOrder, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var updateData model.Order
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve immutable fields
	updateData.ID = id
	updateData.CreatedBy = existingOrder.CreatedBy
	updateData.CreatedAt = existingOrder.CreatedAt
	updateData.BatchID = existingOrder.BatchID

	if err := h.orderRepo.Update(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updateData)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Load existing order to check current status
	existingOrder, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Detect transition to "proses" — trigger auto-deduct stock
	if req.Status == "proses" && existingOrder.ProdStatus != "proses" {
		stockResults, err := h.processStockDeduction(existingOrder)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Status updated successfully, stock deducted",
			"stock_deducted": stockResults,
		})
		return
	}

	if err := h.orderRepo.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *OrderHandler) processStockDeduction(order *model.Order) ([]map[string]interface{}, error) {
	// Aggregate ingredient needs from all order items via recipes
	type ingredientNeed struct {
		ID           uuid.UUID
		Name         string
		Needed       float64
		CurrentStock float64
	}
	needsMap := make(map[uuid.UUID]*ingredientNeed)

	for _, item := range order.Items {
		recipes, err := h.recipeRepo.FindByProductID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to load recipe for product %s: %v", item.ProductName, err)
		}

		for _, recipe := range recipes {
			ingredientID := recipe.IngredientID
			qtyNeeded := recipe.QtyPerBox * float64(item.QtyBox)

			if existing, ok := needsMap[ingredientID]; ok {
				existing.Needed += qtyNeeded
			} else {
				ingredient, err := h.ingredientRepo.FindByID(ingredientID)
				if err != nil {
					return nil, fmt.Errorf("failed to load ingredient %s: %v", ingredientID, err)
				}
				needsMap[ingredientID] = &ingredientNeed{
					ID:           ingredientID,
					Name:         ingredient.Name,
					Needed:       qtyNeeded,
					CurrentStock: ingredient.CurrentStock,
				}
			}
		}
	}

	// Validate stock sufficiency
	var insufficient []string
	for _, need := range needsMap {
		if need.CurrentStock < need.Needed {
			insufficient = append(insufficient, fmt.Sprintf("%s: butuh %.3f, stok %.3f", need.Name, need.Needed, need.CurrentStock))
		}
	}
	if len(insufficient) > 0 {
		return nil, fmt.Errorf("stok tidak cukup: %s", fmt.Sprintf("%v", insufficient))
	}

	// Get user_id from context for stock_log.created_by — but we don't have context here
	// Use order.CreatedBy instead
	userID := order.CreatedBy
	orderID := order.ID
	batchID := order.BatchID

	var results []map[string]interface{}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Update order status
		if err := tx.Model(&model.Order{}).Where("id = ?", orderID).Update("prod_status", "proses").Error; err != nil {
			return err
		}

		for _, need := range needsMap {
			stockBefore := need.CurrentStock
			stockAfter := stockBefore - need.Needed

			// Decrement ingredient stock
			if err := tx.Model(&model.Ingredient{}).Where("id = ?", need.ID).
				Update("current_stock", gorm.Expr("current_stock - ?", need.Needed)).Error; err != nil {
				return err
			}

			// Insert stock_log
			log := &model.StockLog{
				IngredientID:  need.ID,
				BatchID:       batchID,
				LogType:       "out",
				Qty:           need.Needed,
				StockBefore:   stockBefore,
				StockAfter:    stockAfter,
				ReferenceID:   &orderID,
				ReferenceType: "order",
				Note:          fmt.Sprintf("Order #%s - %s", orderID.String()[:8], need.Name),
				CreatedBy:     userID,
			}
			if err := h.stockLogRepo.Create(tx, log); err != nil {
				return err
			}

			results = append(results, map[string]interface{}{
				"ingredient_id":   need.ID,
				"ingredient_name": need.Name,
				"qty_deducted":    need.Needed,
				"stock_before":    stockBefore,
				"stock_after":     stockAfter,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (h *OrderHandler) UpdatePayStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		PayStatus string `json:"pay_status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderRepo.UpdatePayStatus(id, req.PayStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.orderRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
