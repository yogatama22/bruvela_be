package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentStatusHandler struct {
	paymentStatusRepo repository.PaymentStatusRepository
}

func NewPaymentStatusHandler(paymentStatusRepo repository.PaymentStatusRepository) *PaymentStatusHandler {
	return &PaymentStatusHandler{
		paymentStatusRepo: paymentStatusRepo,
	}
}

func (h *PaymentStatusHandler) GetAll(c *gin.Context) {
	paymentStatuses, err := h.paymentStatusRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paymentStatuses)
}

func (h *PaymentStatusHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	paymentStatus, err := h.paymentStatusRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment status not found"})
		return
	}

	c.JSON(http.StatusOK, paymentStatus)
}

func (h *PaymentStatusHandler) Create(c *gin.Context) {
	var paymentStatus model.PaymentStatus
	if err := c.ShouldBindJSON(&paymentStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if status code already exists
	existing, _ := h.paymentStatusRepo.FindByCode(paymentStatus.StatusCode)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Status code already exists"})
		return
	}

	if err := h.paymentStatusRepo.Create(&paymentStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentStatus)
}

func (h *PaymentStatusHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var paymentStatus model.PaymentStatus
	if err := c.ShouldBindJSON(&paymentStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentStatus.ID = id.String()

	// Check if status code already exists (excluding current)
	existing, _ := h.paymentStatusRepo.FindByCode(paymentStatus.StatusCode)
	if existing != nil && existing.ID != id.String() {
		c.JSON(http.StatusConflict, gin.H{"error": "Status code already exists"})
		return
	}

	if err := h.paymentStatusRepo.Update(&paymentStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentStatus)
}

func (h *PaymentStatusHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.paymentStatusRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status deleted successfully"})
}
