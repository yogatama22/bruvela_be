package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderStatusHandler struct {
	orderStatusRepo repository.OrderStatusRepository
}

func NewOrderStatusHandler(orderStatusRepo repository.OrderStatusRepository) *OrderStatusHandler {
	return &OrderStatusHandler{
		orderStatusRepo: orderStatusRepo,
	}
}

func (h *OrderStatusHandler) GetAll(c *gin.Context) {
	orderStatuses, err := h.orderStatusRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderStatuses)
}

func (h *OrderStatusHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	orderStatus, err := h.orderStatusRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order status not found"})
		return
	}

	c.JSON(http.StatusOK, orderStatus)
}

func (h *OrderStatusHandler) Create(c *gin.Context) {
	var orderStatus model.OrderStatus
	if err := c.ShouldBindJSON(&orderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if status code already exists
	existing, _ := h.orderStatusRepo.FindByCode(orderStatus.OrderStatusCode)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Status code already exists"})
		return
	}

	if err := h.orderStatusRepo.Create(&orderStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderStatus)
}

func (h *OrderStatusHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var orderStatus model.OrderStatus
	if err := c.ShouldBindJSON(&orderStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderStatus.ID = id.String()

	// Check if status code already exists (excluding current)
	existing, _ := h.orderStatusRepo.FindByCode(orderStatus.OrderStatusCode)
	if existing != nil && existing.ID != id.String() {
		c.JSON(http.StatusConflict, gin.H{"error": "Status code already exists"})
		return
	}

	if err := h.orderStatusRepo.Update(&orderStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderStatus)
}

func (h *OrderStatusHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.orderStatusRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status deleted successfully"})
}
