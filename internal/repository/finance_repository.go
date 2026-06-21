package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceRepository interface {
	CreateJournalEntry(entry *model.JournalEntry) error
	FindJournalEntries(batchID uuid.UUID) ([]model.JournalEntry, error)
	GetBatchSummary(batchID uuid.UUID) (map[string]interface{}, error)
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) CreateJournalEntry(entry *model.JournalEntry) error {
	return r.db.Create(entry).Error
}

func (r *financeRepository) FindJournalEntries(batchID uuid.UUID) ([]model.JournalEntry, error) {
	var entries []model.JournalEntry
	err := r.db.Where("batch_id = ?", batchID).Order("entry_date DESC").Find(&entries).Error
	return entries, err
}

func (r *financeRepository) GetBatchSummary(batchID uuid.UUID) (map[string]interface{}, error) {
	var result struct {
		TotalIncome  int
		TotalExpense int
		Balance      int
	}

	err := r.db.Model(&model.JournalEntry{}).
		Select("SUM(CASE WHEN type IN ('income', 'modal') THEN amount ELSE 0 END) as total_income, "+
			"SUM(CASE WHEN type IN ('expense', 'transfer') THEN amount ELSE 0 END) as total_expense, "+
			"SUM(amount) as balance").
		Where("batch_id = ?", batchID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_income":  result.TotalIncome,
		"total_expense": result.TotalExpense,
		"balance":       result.Balance,
	}, nil
}
