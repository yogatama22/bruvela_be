package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShippingTypeHandler struct {
	shippingTypeRepo repository.ShippingTypeRepository
}

func NewShippingTypeHandler(shippingTypeRepo repository.ShippingTypeRepository) *ShippingTypeHandler {
	return &ShippingTypeHandler{
		shippingTypeRepo: shippingTypeRepo,
	}
}

func (h *ShippingTypeHandler) GetAll(c *gin.Context) {
	shippingTypes, err := h.shippingTypeRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shippingTypes)
}

func (h *ShippingTypeHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	shippingType, err := h.shippingTypeRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping type not found"})
		return
	}

	c.JSON(http.StatusOK, shippingType)
}

func (h *ShippingTypeHandler) Create(c *gin.Context) {
	var shippingType model.ShippingType
	if err := c.ShouldBindJSON(&shippingType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if shipping code already exists
	existing, _ := h.shippingTypeRepo.FindByCode(shippingType.ShippingCode)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Shipping code already exists"})
		return
	}

	if err := h.shippingTypeRepo.Create(&shippingType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shippingType)
}

func (h *ShippingTypeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var shippingType model.ShippingType
	if err := c.ShouldBindJSON(&shippingType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shippingType.ID = id.String()

	// Check if shipping code already exists (excluding current)
	existing, _ := h.shippingTypeRepo.FindByCode(shippingType.ShippingCode)
	if existing != nil && existing.ID != id.String() {
		c.JSON(http.StatusConflict, gin.H{"error": "Shipping code already exists"})
		return
	}

	if err := h.shippingTypeRepo.Update(&shippingType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shippingType)
}

func (h *ShippingTypeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.shippingTypeRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping type deleted successfully"})
}
