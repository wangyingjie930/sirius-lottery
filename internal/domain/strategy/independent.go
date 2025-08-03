package strategy

import (
	"context"
	"errors"
	"math/rand"
	"sirius-lottery/internal/domain"
	"sirius-lottery/internal/domain/entity"
	"sort"
)

// --- 策略1: 独立概率 (Independent Probability) ---

type independentProbabilityStrategy struct{}

// NewIndependentProbabilityStrategy 创建一个独立概率策略实例
func NewIndependentProbabilityStrategy() domain.LotteryStrategy {
	return &independentProbabilityStrategy{}
}

func (s *independentProbabilityStrategy) Draw(ctx context.Context, drawCtx *domain.DrawContext) (*entity.LotteryPrize, error) {
	// 按概率从高到低排序，优先判断高概率奖品
	sort.Slice(drawCtx.Prizes, func(i, j int) bool {
		return drawCtx.Prizes[i].Probability > drawCtx.Prizes[j].Probability
	})

	// 生成一个 [0.0, 1.0) 的随机数
	random := rand.Float64()
	var currentProb float64 = 0

	for _, prize := range drawCtx.Prizes {
		// IsSpecial 通常指 "谢谢参与"，它不参与概率计算，作为默认奖品
		if prize.IsSpecial {
			continue
		}
		currentProb += prize.Probability
		if random < currentProb {
			return prize, nil // 命中奖品
		}
	}

	// 如果所有奖品都未命中，则返回特殊奖品（如谢谢参与）
	for _, prize := range drawCtx.Prizes {
		if prize.IsSpecial {
			return prize, nil
		}
	}

	return nil, errors.New("抽奖失败：奖池中未配置默认奖品 (e.g., 谢谢参与)")
}
