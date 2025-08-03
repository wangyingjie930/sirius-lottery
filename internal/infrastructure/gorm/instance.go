package gorm

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// LotteryInstance 抽奖实例表
type LotteryInstance struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	InstanceID    string    `gorm:"type:varchar(50);uniqueIndex:uk_instance_id;not null;comment:业务活动ID" json:"instance_id"`
	InstanceName  string    `gorm:"type:varchar(255);not null;comment:活动名称" json:"instance_name"`
	TemplateID    uint64    `gorm:"not null;comment:关联的模板ID" json:"template_id"`
	StartTime     time.Time `gorm:"type:timestamp;not null;index:idx_start_end_time;comment:活动开始时间" json:"start_time"`
	EndTime       time.Time `gorm:"type:timestamp;not null;index:idx_start_end_time;comment:活动结束时间" json:"end_time"`
	UserScopeJSON JSONMap   `gorm:"type:json;comment:参与用户限制" json:"user_scope_json"`
	Status        int8      `gorm:"type:tinyint;not null;default:1;comment:状态: 1-待上线, 2-进行中, 3-已下线" json:"status"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Template *LotteryTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Pools    []LotteryPool    `gorm:"foreignKey:InstanceID;references:InstanceID" json:"pools,omitempty"`
}

// TableName 指定表名
func (LotteryInstance) TableName() string {
	return "lottery_instance"
}

func (li *LotteryInstance) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑，比如验证时间范围
	if li.EndTime.Before(li.StartTime) {
		return errors.New("结束时间不能早于开始时间")
	}
	return nil
}
