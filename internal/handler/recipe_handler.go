package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler struct {
	recipeRepo     repository.RecipeRepository
	productRepo    repository.ProductRepository
	ingredientRepo repository.IngredientRepository
}

func NewRecipeHandler(recipeRepo repository.RecipeRepository, productRepo repository.ProductRepository, ingredientRepo repository.IngredientRepository) *RecipeHandler {
	return &RecipeHandler{
		recipeRepo:     recipeRepo,
		productRepo:    productRepo,
		ingredientRepo: ingredientRepo,
	}
}

type ProductWithHPP struct {
	ID        uuid.UUID      `json:"id"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	Price     int            `json:"price"`
	PcsPerBox int            `json:"pcs_per_box"`
	Status    string         `json:"status"`
	HPP       int            `json:"hpp"`
	Margin    float64        `json:"margin"`
	Profit    int            `json:"profit"`
	Recipes   []model.Recipe `json:"recipes,omitempty"`
}

func (h *RecipeHandler) GetAllProductsWithHPP(c *gin.Context) {
	products, err := h.productRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var productsWithHPP []ProductWithHPP

	for _, product := range products {
		recipes, err := h.recipeRepo.FindByProductID(product.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		hpp := 0.0
		for _, recipe := range recipes {
			hpp += recipe.CostPerBox
		}

		profit := float64(product.Price) - hpp
		margin := 0.0
		if product.Price > 0 {
			margin = (profit / float64(product.Price)) * 100
		}

		productsWithHPP = append(productsWithHPP, ProductWithHPP{
			ID:        product.ID,
			Code:      product.Code,
			Name:      product.Name,
			Price:     product.Price,
			PcsPerBox: product.PcsPerBox,
			Status:    product.Status,
			HPP:       int(hpp),
			Margin:    margin,
			Profit:    int(profit),
			Recipes:   recipes,
		})
	}

	c.JSON(http.StatusOK, productsWithHPP)
}

func (h *RecipeHandler) GetByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	recipes, err := h.recipeRepo.FindByProductID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func (h *RecipeHandler) Create(c *gin.Context) {
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch ingredient to get price_per_use
	ingredient, err := h.ingredientRepo.FindByID(recipe.IngredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ingredient"})
		return
	}

	// Calculate cost_per_box
	recipe.CostPerBox = recipe.QtyPerBox * ingredient.PricePerUse

	if err := h.recipeRepo.Create(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

func (h *RecipeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = id

	// Fetch ingredient to get price_per_use
	ingredient, err := h.ingredientRepo.FindByID(recipe.IngredientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ingredient"})
		return
	}

	// Calculate cost_per_box
	recipe.CostPerBox = recipe.QtyPerBox * ingredient.PricePerUse

	if err := h.recipeRepo.Update(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.recipeRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

func (h *RecipeHandler) DeleteByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.recipeRepo.DeleteByProductID(productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipes deleted successfully"})
}

type CalculatorItem struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	QtyBox    int       `json:"qty_box" binding:"required,gt=0"`
}

type CalculatorRequest struct {
	Items []CalculatorItem `json:"items" binding:"required"`
}

type CalculatorIngredient struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Needed       float64   `json:"needed"`
	CurrentStock float64   `json:"current_stock"`
	Status       string    `json:"status"`
}

type CalculatorResponse struct {
	Ingredients []CalculatorIngredient `json:"ingredients"`
	TotalHPP    float64                `json:"total_hpp"`
}

func (h *RecipeHandler) CalculateProduction(c *gin.Context) {
	var req CalculatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type ingredientAgg struct {
		ID           uuid.UUID
		Name         string
		Needed       float64
		CurrentStock float64
	}
	aggMap := make(map[uuid.UUID]*ingredientAgg)

	for _, item := range req.Items {
		recipes, err := h.recipeRepo.FindByProductID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load recipe"})
			return
		}

		for _, recipe := range recipes {
			qtyNeeded := recipe.QtyPerBox * float64(item.QtyBox)
			if existing, ok := aggMap[recipe.IngredientID]; ok {
				existing.Needed += qtyNeeded
			} else {
				ingredient, err := h.ingredientRepo.FindByID(recipe.IngredientID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ingredient"})
					return
				}
				aggMap[recipe.IngredientID] = &ingredientAgg{
					ID:           recipe.IngredientID,
					Name:         ingredient.Name,
					Needed:       qtyNeeded,
					CurrentStock: ingredient.CurrentStock,
				}
			}
		}
	}

	var ingredients []CalculatorIngredient
	totalHPP := 0.0

	for _, agg := range aggMap {
		status := "cukup"
		if agg.CurrentStock < agg.Needed {
			status = "kurang"
		} else if agg.CurrentStock < agg.Needed*1.1 {
			status = "pas-pasan"
		}

		ingredients = append(ingredients, CalculatorIngredient{
			ID:           agg.ID,
			Name:         agg.Name,
			Needed:       agg.Needed,
			CurrentStock: agg.CurrentStock,
			Status:       status,
		})

		ingredient, _ := h.ingredientRepo.FindByID(agg.ID)
		if ingredient != nil {
			totalHPP += agg.Needed * ingredient.PricePerUse
		}
	}

	c.JSON(http.StatusOK, CalculatorResponse{
		Ingredients: ingredients,
		TotalHPP:    totalHPP,
	})
}
