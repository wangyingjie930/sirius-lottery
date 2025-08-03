package entity

const (
	WinRecordStatusPending = 1 // 待发放
	WinRecordStatusSuccess = 2 // 发放成功
	WinRecordStatusFailed  = 3 // 发放失败
)

type LotteryWinRecord struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID  string
	OrderID    string `gorm:"type:varchar(64);uniqueIndex:uk_order_id;not null;comment:唯一订单号" json:"order_id"`
	InstanceID string `gorm:"type:varchar(50);index:idx_user_instance;not null" json:"instance_id"`
	UserID     uint64 `gorm:"index:idx_user_instance;not null" json:"user_id"`
	PrizeID    string `gorm:"type:varchar(100);not null" json:"prize_id"`
	Status     int8   `gorm:"type:tinyint;not null;default:1;comment:发放状态: 1-待发放, 2-发放成功, 3-发放失败" json:"status"`
}

func (l *LotteryWinRecord) IsThankYouPrize() bool {
	return false
}
