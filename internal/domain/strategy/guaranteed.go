package strategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sirius-lottery/internal/domain"
	"sirius-lottery/internal/domain/entity"
)

// GuaranteeRepository 定义了保底计数器所需的数据操作
// 其具体实现应该在 infrastructure 层，例如使用 Redis
type GuaranteeRepository interface {
	// IncrementAndGet 返回自增后的连续未中奖次数
	IncrementAndGet(ctx context.Context, instanceID string, userID int64) (int, error)
	// ResetCounter 重置用户的未中奖次数
	ResetCounter(ctx context.Context, instanceID string, userID int64) error
}

type guaranteedWinStrategy struct {
	guaranteeRepo GuaranteeRepository
}

// NewGuaranteedWinStrategy 创建一个保底策略实例
func NewGuaranteedWinStrategy(repo GuaranteeRepository) domain.LotteryStrategy {
	return &guaranteedWinStrategy{guaranteeRepo: repo}
}

// 保底策略的配置结构
type guaranteeConfig struct {
	GuaranteeCount int `json:"guarantee_count"` // 需要保底的次数
}

func (s *guaranteedWinStrategy) Draw(ctx context.Context, drawCtx *domain.DrawContext) (*entity.LotteryPrize, error) {
	// 1. 解析策略配置
	var config guaranteeConfig
	if err := json.Unmarshal([]byte(drawCtx.Pool.StrategyConfigJSON), &config); err != nil {
		return nil, errors.New("保底策略配置解析失败")
	}
	if config.GuaranteeCount <= 0 {
		return nil, errors.New("保底次数必须大于0")
	}

	// 2. 获取用户连续未中奖次数
	missCount, err := s.guaranteeRepo.IncrementAndGet(ctx, drawCtx.InstanceID, drawCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取保底计数失败: %w", err)
	}

	// 3. 判断是否达到保底次数
	if missCount >= config.GuaranteeCount {
		// 达到保底，从非特殊奖品中随机抽取一个
		var normalPrizes []*entity.LotteryPrize
		for _, p := range drawCtx.Prizes {
			if !p.IsSpecial {
				normalPrizes = append(normalPrizes, p)
			}
		}
		if len(normalPrizes) == 0 {
			return nil, errors.New("保底触发失败：奖池中没有可保底的普通奖品")
		}

		// 随机选择一个保底奖品
		wonPrize := normalPrizes[rand.Intn(len(normalPrizes))]

		// 中奖后重置计数器
		_ = s.guaranteeRepo.ResetCounter(ctx, drawCtx.InstanceID, drawCtx.UserID)
		return wonPrize, nil

	} else {
		// 4. 未达到保底次数，走独立概率逻辑
		// 这里可以直接复用上面的 independentProbabilityStrategy 逻辑
		fallbackStrategy := NewIndependentProbabilityStrategy()
		wonPrize, err := fallbackStrategy.Draw(ctx, drawCtx)
		if err != nil {
			return nil, err
		}

		// 如果中的不是特殊奖品（即中奖了），则重置计数器
		if !wonPrize.IsSpecial {
			_ = s.guaranteeRepo.ResetCounter(ctx, drawCtx.InstanceID, drawCtx.UserID)
		}

		return wonPrize, nil
	}
}
