package application

import (
	"context"
	"sirius-lottery/internal/pkg/eventbus"
)

type LotteryService interface {
	// Draw 是核心抽奖接口
	Draw(ctx context.Context, req *DrawRequest) (*DrawResponse, error)

	// GetLotteryInstance 获取活动详情，用于前端渲染 [cite: 167]
	GetLotteryInstance(ctx context.Context, instanceID string) (*LotteryInstanceResponse, error)

	DeductStock(ctx context.Context, req *StockActionRequest) error
	IncreaseStock(ctx context.Context, req *StockActionRequest) error
}

type EventBus interface {
	Publish(ctx context.Context, event eventbus.Event) error
}
