package gorm

import "time"

// LotteryPrize 奖品配置表
type LotteryPrize struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PoolID         uint64    `gorm:"index:idx_pool_id;not null;comment:关联的奖池ID" json:"pool_id"`
	PrizeID        string    `gorm:"type:varchar(100);not null;comment:业务奖品ID" json:"prize_id"`
	PrizeName      string    `gorm:"type:varchar(255);not null;comment:奖品名" json:"prize_name"`
	AllocatedStock int       `gorm:"not null;default:0;comment:总预算库存" json:"allocated_stock"`
	Probability    float64   `gorm:"type:decimal(10,8);not null;default:0.00000000;comment:中奖概率" json:"probability"`
	IsSpecial      bool      `gorm:"type:tinyint(1);not null;default:0;comment:是否特殊奖品" json:"is_special"`
	CreatedAt      time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Pool *LotteryPool `gorm:"foreignKey:PoolID" json:"pool,omitempty"`
}

// TableName 指定表名
func (LotteryPrize) TableName() string {
	return "lottery_prize"
}
