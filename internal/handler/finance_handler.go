package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FinanceHandler struct {
	financeRepo repository.FinanceRepository
}

func NewFinanceHandler(financeRepo repository.FinanceRepository) *FinanceHandler {
	return &FinanceHandler{financeRepo: financeRepo}
}

func (h *FinanceHandler) CreateJournalEntry(c *gin.Context) {
	var entry model.JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	entry.CreatedBy = userID.(uuid.UUID)

	if err := h.financeRepo.CreateJournalEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *FinanceHandler) GetJournalEntries(c *gin.Context) {
	batchID, err := uuid.Parse(c.Query("batch_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch_id"})
		return
	}

	entries, err := h.financeRepo.FindJournalEntries(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (h *FinanceHandler) GetSummary(c *gin.Context) {
	batchID, err := uuid.Parse(c.Query("batch_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch_id"})
		return
	}

	summary, err := h.financeRepo.GetBatchSummary(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
