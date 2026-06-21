package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseHandler struct {
	purchaseRepo   repository.PurchaseRepository
	ingredientRepo repository.IngredientRepository
	journalRepo    repository.FinanceRepository
	stockLogRepo   repository.StockLogRepository
	db             *gorm.DB
}

func NewPurchaseHandler(
	purchaseRepo repository.PurchaseRepository,
	ingredientRepo repository.IngredientRepository,
	journalRepo repository.FinanceRepository,
	stockLogRepo repository.StockLogRepository,
	db *gorm.DB,
) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseRepo:   purchaseRepo,
		ingredientRepo: ingredientRepo,
		journalRepo:    journalRepo,
		stockLogRepo:   stockLogRepo,
		db:             db,
	}
}

func (h *PurchaseHandler) GetAll(c *gin.Context) {
	filters := make(map[string]interface{})
	if batchID := c.Query("batch_id"); batchID != "" {
		filters["batch_id"] = batchID
	}
	if ingredientID := c.Query("ingredient_id"); ingredientID != "" {
		filters["ingredient_id"] = ingredientID
	}

	purchases, err := h.purchaseRepo.FindAll(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchases)
}

func (h *PurchaseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	purchase, err := h.purchaseRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase not found"})
		return
	}
	c.JSON(http.StatusOK, purchase)
}

// Create creates a purchase and:
//  1. Adds stock to ingredient (current_stock += qty_pack * qty_per_pack)
//  2. Auto-creates a journal entry (type=expense) for accounting
//
// All in a single DB transaction so partial failure rolls back atomically.
func (h *PurchaseHandler) Create(c *gin.Context) {
	var purchase model.IngredientPurchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	if uid, ok := userID.(uuid.UUID); ok {
		purchase.CreatedBy = uid
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create purchase record (BeforeCreate hook computes TotalPrice)
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		// 2. Fetch ingredient to know qty_per_pack
		var ingredient model.Ingredient
		if err := tx.First(&ingredient, "id = ?", purchase.IngredientID).Error; err != nil {
			return err
		}

		// 3. Add stock: qty_pack * qty_per_pack = total use-units to add
		stockToAdd := purchase.QtyPack * ingredient.QtyPerPack
		if err := tx.Model(&model.Ingredient{}).
			Where("id = ?", purchase.IngredientID).
			Update("current_stock", gorm.Expr("current_stock + ?", stockToAdd)).Error; err != nil {
			return err
		}

		// 4. Auto-create journal entry (outflow = negative)
		supplierLabel := purchase.Supplier
		if supplierLabel == "" {
			supplierLabel = "Supplier"
		}
		journalEntry := model.JournalEntry{
			BatchID:     purchase.BatchID,
			EntryDate:   purchase.PurchaseDate.Time,
			Description: "Pembelian " + ingredient.Name + " - " + supplierLabel,
			Type:        "expense",
			Amount:      -purchase.TotalPrice,
			CreatedBy:   purchase.CreatedBy,
		}
		if err := tx.Create(&journalEntry).Error; err != nil {
			return err
		}

		// 5. Insert stock_log (log_type="in", reference_type="purchase")
		purchaseID := purchase.ID
		stockLog := &model.StockLog{
			IngredientID:  purchase.IngredientID,
			BatchID:       purchase.BatchID,
			LogType:       "in",
			Qty:           stockToAdd,
			StockBefore:   ingredient.CurrentStock,
			StockAfter:    ingredient.CurrentStock + stockToAdd,
			ReferenceID:   &purchaseID,
			ReferenceType: "purchase",
			Note:          "Pembelian " + ingredient.Name + " - " + supplierLabel,
			CreatedBy:     purchase.CreatedBy,
		}
		if err := h.stockLogRepo.Create(tx, stockLog); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return with preloaded relations
	if loaded, lerr := h.purchaseRepo.FindByID(purchase.ID); lerr == nil {
		c.JSON(http.StatusCreated, loaded)
		return
	}
	c.JSON(http.StatusCreated, purchase)
}

func (h *PurchaseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var purchase model.IngredientPurchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	purchase.ID = id

	if err := h.purchaseRepo.Update(&purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.purchaseRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted successfully"})
}
