package entity

type LotteryPool struct {
	ID                 int64
	InstanceID         string
	PoolName           string
	CostJSON           string // 消耗的资产列表, e.g., [{"asset_id": "ticket", "amount": 1}] [cite: 140]
	LotteryStrategy    string // 抽奖算法策略 [cite: 140]
	StrategyConfigJSON string // 策略相关配置, e.g., {"guarantee_count": 10} [cite: 140]

	Prizes []*LotteryPrize
}

func (l *LotteryPool) GetCost() int {
	return 0
}
