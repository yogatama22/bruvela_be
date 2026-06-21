package handler

import (
	"bruvela-backend/internal/model"
	"bruvela-backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompanySettingHandler struct {
	settingRepo repository.CompanySettingRepository
}

func NewCompanySettingHandler(settingRepo repository.CompanySettingRepository) *CompanySettingHandler {
	return &CompanySettingHandler{settingRepo: settingRepo}
}

func (h *CompanySettingHandler) Get(c *gin.Context) {
	setting, err := h.settingRepo.Get()
	if err != nil {
		c.JSON(http.StatusOK, &model.CompanySetting{})
		return
	}
	c.JSON(http.StatusOK, setting)
}

func (h *CompanySettingHandler) Upsert(c *gin.Context) {
	var setting model.CompanySetting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingRepo.Upsert(&setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setting)
}
