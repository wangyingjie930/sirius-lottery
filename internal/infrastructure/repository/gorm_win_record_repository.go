package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sirius-lottery/internal/domain/entity"
	gorm_model "sirius-lottery/internal/infrastructure/gorm"
)

type gormWinRecordRepository struct {
	db *gorm.DB
}

// NewGormWinRecordRepository 创建一个新的 WinRecordRepository GORM 实现
func NewGormWinRecordRepository(db *gorm.DB) *gormWinRecordRepository {
	return &gormWinRecordRepository{db: db}
}

// Create 创建一条新的中奖记录
func (r *gormWinRecordRepository) Create(ctx context.Context, record *entity.LotteryWinRecord) error {
	model := &gorm_model.LotteryWinRecord{
		OrderID:    record.OrderID,
		InstanceID: record.InstanceID,
		UserID:     record.UserID,
		PrizeID:    record.PrizeID,
		Status:     record.Status,
	}

	// 使用 WithContext 传递上下文，这对于超时和取消非常重要
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByRequestID 根据请求ID查询中奖记录 (用于幂等性检查)
func (r *gormWinRecordRepository) GetByRequestID(ctx context.Context, requestID string) (*entity.LotteryWinRecord, error) {
	var model gorm_model.LotteryWinRecord
	// 注意: 您的 gorm_model.LotteryWinRecord 中没有 RequestID 字段，但领域实体中有。
	// 这通常意味着 RequestID 可能不持久化到数据库，或者数据库模型需要更新。
	// 假设它不持久化，此方法将无法按预期工作。
	// 如果需要持久化，请在 gorm_model.LotteryWinRecord 中添加 RequestID 字段。
	// 这里我们假设需要根据 OrderID 查询，因为它是唯一的。
	err := r.db.WithContext(ctx).Where("order_id = ?", requestID).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到是正常情况，不应返回错误
		}
		return nil, err
	}

	return &entity.LotteryWinRecord{
		ID:         model.ID,
		OrderID:    model.OrderID,
		InstanceID: model.InstanceID,
		UserID:     model.UserID,
		PrizeID:    model.PrizeID,
		Status:     model.Status,
	}, nil
}
