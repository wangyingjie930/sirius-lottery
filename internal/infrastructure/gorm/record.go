package gorm

import "time"

// LotteryWinRecord 中奖记录表
type LotteryWinRecord struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID    string    `gorm:"type:varchar(64);uniqueIndex:uk_order_id;not null;comment:唯一订单号" json:"order_id"`
	InstanceID string    `gorm:"type:varchar(50);index:idx_user_instance;not null" json:"instance_id"`
	UserID     uint64    `gorm:"index:idx_user_instance;not null" json:"user_id"`
	PrizeID    string    `gorm:"type:varchar(100);not null" json:"prize_id"`
	Status     int8      `gorm:"type:tinyint;not null;default:1;comment:发放状态: 1-待发放, 2-发放成功, 3-发放失败" json:"status"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Instance *LotteryInstance `gorm:"foreignKey:InstanceID;references:InstanceID" json:"instance,omitempty"`
}

// TableName 指定表名
func (LotteryWinRecord) TableName() string {
	return "lottery_win_record"
}
