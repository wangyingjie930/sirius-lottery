package gorm

import "time"

// LotteryPool 奖池配置表
type LotteryPool struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	InstanceID         string    `gorm:"type:varchar(50);index:idx_instance_id;not null;comment:关联的抽奖实例ID" json:"instance_id"`
	PoolName           string    `gorm:"type:varchar(100);not null;comment:奖池名称" json:"pool_name"`
	CostJSON           JSONArray `gorm:"type:json;not null;comment:消耗的资产列表" json:"cost_json"`
	LotteryStrategy    string    `gorm:"type:varchar(50);not null;comment:抽奖算法策略" json:"lottery_strategy"`
	StrategyConfigJSON JSONMap   `gorm:"type:json;comment:策略相关配置" json:"strategy_config_json"`
	CreatedAt          time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Instance *LotteryInstance `gorm:"foreignKey:InstanceID;references:InstanceID" json:"instance,omitempty"`
	Prizes   []LotteryPrize   `gorm:"foreignKey:PoolID" json:"prizes,omitempty"`
}

// TableName 指定表名
func (LotteryPool) TableName() string {
	return "lottery_pool"
}
