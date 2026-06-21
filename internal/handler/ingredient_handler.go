package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IngredientHandler struct {
	ingredientRepo repository.IngredientRepository
	recipeRepo     repository.RecipeRepository
	orderRepo      repository.OrderRepository
}

func NewIngredientHandler(ingredientRepo repository.IngredientRepository, recipeRepo repository.RecipeRepository, orderRepo repository.OrderRepository) *IngredientHandler {
	return &IngredientHandler{
		ingredientRepo: ingredientRepo,
		recipeRepo:     recipeRepo,
		orderRepo:      orderRepo,
	}
}

func (h *IngredientHandler) GetAll(c *gin.Context) {
	ingredients, err := h.ingredientRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

func (h *IngredientHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ingredient, err := h.ingredientRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func (h *IngredientHandler) Create(c *gin.Context) {
	var ingredient model.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ingredientRepo.Create(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ingredient)
}

func (h *IngredientHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ingredient model.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredient.ID = id
	if err := h.ingredientRepo.Update(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.ingredientRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted successfully"})
}

func (h *IngredientHandler) GetLowStock(c *gin.Context) {
	ingredients, err := h.ingredientRepo.FindLowStock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

type IngredientWithEstimation struct {
	model.Ingredient
	EstimatedUsage float64 `json:"estimated_usage"`
}

func (h *IngredientHandler) GetWithEstimation(c *gin.Context) {
	batchID := c.Query("batch_id")

	ingredients, err := h.ingredientRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pre-calculate estimation per ingredient if batch is selected
	estimationMap := make(map[uuid.UUID]float64)

	if batchID != "" {
		batchUUID, err := uuid.Parse(batchID)
		if err == nil {
			// Get orders for this batch with items preloaded
			orders, err := h.orderRepo.FindByBatchID(batchUUID)
			if err == nil {
				log.Printf("Found %d orders for batch %s", len(orders), batchID)

				for _, order := range orders {
					log.Printf("Order %s has %d items", order.ID, len(order.Items))

					for _, item := range order.Items {
						// Get recipes for this product
						recipes, err := h.recipeRepo.FindByProductID(item.ProductID)
						if err != nil {
							continue
						}
						log.Printf("Product %s (%s) has %d recipes, qty_box: %d", item.ProductID, item.ProductName, len(recipes), item.QtyBox)

						for _, recipe := range recipes {
							qty := recipe.QtyPerBox * float64(item.QtyBox)
							estimationMap[recipe.IngredientID] += qty
							log.Printf("  Recipe ingredient %s: %.2f x %d = %.2f", recipe.IngredientID, recipe.QtyPerBox, item.QtyBox, qty)
						}
					}
				}
			} else {
				log.Printf("Error fetching orders: %v", err)
			}
		}
	}

	var result []IngredientWithEstimation
	for _, ingredient := range ingredients {
		result = append(result, IngredientWithEstimation{
			Ingredient:     ingredient,
			EstimatedUsage: estimationMap[ingredient.ID],
		})
	}

	c.JSON(http.StatusOK, result)
}
