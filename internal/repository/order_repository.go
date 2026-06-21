package repository

import (
	"bruvela-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id uuid.UUID) (*model.Order, error)
	FindAll(filters map[string]interface{}, limit, offset int) ([]model.Order, int64, error)
	FindByBatchID(batchID uuid.UUID) ([]model.Order, error)
	Update(order *model.Order) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status string) error
	UpdatePayStatus(id uuid.UUID, payStatus string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	// Create order with items in a transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Create the order first
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// After items are created with their subtotals, recalculate order totals
		order.CalculateTotals()

		// Update the order with calculated totals
		if err := tx.Model(order).Updates(map[string]interface{}{
			"total_product": order.TotalProduct,
			"total_bill":    order.TotalBill,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *orderRepository) FindByID(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items.Product").Preload("Batch").Preload("Customer").First(&order, "id = ?", id).Error
	return &order, err
}

func (r *orderRepository) FindAll(filters map[string]interface{}, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Items").Preload("Batch").Preload("Customer").
		Limit(limit).Offset(offset).Order("order_date DESC").Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) FindByBatchID(batchID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items").Preload("Items.Product").
		Where("batch_id = ?", batchID).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the order
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		// Delete existing items for this order
		if err := tx.Where("order_id = ?", order.ID).Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}

		// Create new items
		if len(order.Items) > 0 {
			for i := range order.Items {
				order.Items[i].OrderID = order.ID
				if err := tx.Create(&order.Items[i]).Error; err != nil {
					return err
				}
			}
		}

		// Recalculate totals
		order.CalculateTotals()

		// Update the order with calculated totals
		if err := tx.Model(order).Updates(map[string]interface{}{
			"total_product": order.TotalProduct,
			"total_bill":    order.TotalBill,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Order{}, "id = ?", id).Error
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("prod_status", status).Error
}

func (r *orderRepository) UpdatePayStatus(id uuid.UUID, payStatus string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("pay_status", payStatus).Error
}
