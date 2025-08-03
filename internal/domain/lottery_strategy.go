package domain

import (
	"context"
	"sirius-lottery/internal/domain/entity"
)

const (
	StrategyIndependentProbability = "independent_probability" // 独立概率
	StrategyGuaranteedWin          = "guaranteed_win"          // 保底中奖
)

// DrawContext 包含一次抽奖所需的所有上下文信息 [cite: 93]
type DrawContext struct {
	InstanceID string
	UserID     int64
	Pool       *entity.LotteryPool
	Prizes     []*entity.LotteryPrize
}

// LotteryStrategy 定义了抽奖算法的统一接口 [cite: 93]
type LotteryStrategy interface {
	// Draw 执行抽奖算法，返回中奖的奖品
	Draw(ctx context.Context, drawCtx *DrawContext) (*entity.LotteryPrize, error)
}
