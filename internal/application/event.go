package application

import "sirius-lottery/internal/domain/entity"

const (
	TopicLotteryWin = "lottery_win_events"
)

type LotteryWinEvent struct {
	WinRecord *entity.LotteryWinRecord
}

func (e *LotteryWinEvent) Topic() string {
	return TopicLotteryWin
}
