package strategy

import (
	"fmt"
	"sirius-lottery/internal/domain"
	"sync"
)

// LotteryStrategyFactory 是一个用于创建和管理抽奖策略的工厂
// 使用单例模式确保全局只有一个工厂实例
type LotteryStrategyFactory struct {
	strategies map[string]domain.LotteryStrategy
}

var (
	factoryInstance *LotteryStrategyFactory
	once            sync.Once
)

// NewLotteryStrategyFactory 创建并返回一个策略工厂的单例
// 在这里注入所有策略所需的依赖，例如 Redis 客户端
func NewLotteryStrategyFactory(guaranteeRepo GuaranteeRepository) *LotteryStrategyFactory {
	once.Do(func() {
		strategies := make(map[string]domain.LotteryStrategy)
		// 注册所有可用的策略
		strategies[domain.StrategyIndependentProbability] = NewIndependentProbabilityStrategy()
		strategies[domain.StrategyGuaranteedWin] = NewGuaranteedWinStrategy(guaranteeRepo)

		factoryInstance = &LotteryStrategyFactory{
			strategies: strategies,
		}
	})
	return factoryInstance
}

// GetStrategy 根据策略名称返回一个具体的策略实例
func (f *LotteryStrategyFactory) GetStrategy(strategyName string) (domain.LotteryStrategy, error) {
	strategy, ok := f.strategies[strategyName]
	if !ok {
		return nil, fmt.Errorf("未找到名为 '%s' 的抽奖策略", strategyName)
	}
	return strategy, nil
}
