package gorm

import (
	"gorm.io/gorm"
	"time"
)

const (
	TemplateStatusDraft     = 1 // 草稿
	TemplateStatusPublished = 2 // 已发布
	TemplateStatusArchived  = 3 // 已归档

	UIStyleLuckyWheel = "lucky_wheel" // 幸运大转盘
	UIStyleNineGrid   = "nine_grid"   // 九宫格
)

// LotteryTemplate 抽奖模板表
type LotteryTemplate struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	TemplateName string    `gorm:"type:varchar(100);not null;comment:模板名称" json:"template_name"`
	UIStyle      string    `gorm:"type:varchar(50);not null;comment:UI样式标识" json:"ui_style"`
	ConfigJSON   JSONMap   `gorm:"type:json;comment:UI相关的配置" json:"config_json"`
	Status       int8      `gorm:"type:tinyint;not null;default:1;comment:状态: 1-草稿, 2-已发布, 3-已归档" json:"status"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (LotteryTemplate) TableName() string {
	return "lottery_template"
}

func (lt *LotteryTemplate) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}
