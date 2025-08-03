package entity

import (
	"time"
)

const (
	InstanceStatusPending = 1 // 待上线
	InstanceStatusActive  = 2 // 进行中
	InstanceStatusOffline = 3 // 已下线
)

type LotteryInstance struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	InstanceID   string    `gorm:"type:varchar(50);uniqueIndex:uk_instance_id;not null;comment:业务活动ID" json:"instance_id"`
	InstanceName string    `gorm:"type:varchar(255);not null;comment:活动名称" json:"instance_name"`
	TemplateID   uint64    `gorm:"not null;comment:关联的模板ID" json:"template_id"`
	StartTime    time.Time `gorm:"type:timestamp;not null;index:idx_start_end_time;comment:活动开始时间" json:"start_time"`
	EndTime      time.Time `gorm:"type:timestamp;not null;index:idx_start_end_time;comment:活动结束时间" json:"end_time"`
	Status       int8      `gorm:"type:tinyint;not null;default:1;comment:状态: 1-待上线, 2-进行中, 3-已下线" json:"status"`

	Pools []LotteryPool `gorm:"foreignKey:InstanceID;references:InstanceID" json:"pools,omitempty"`
}

func (l *LotteryInstance) Check(cur time.Time) error {
	return nil
}

func (l *LotteryInstance) IsUserAllowed(userId int64) bool {
	return true
}
