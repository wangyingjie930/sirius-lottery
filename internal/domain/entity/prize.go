package entity

// LotteryPrize 定义了奖池中的奖品 [cite: 141]
type LotteryPrize struct {
	ID             int64
	PoolID         int64
	PrizeID        string // 业务奖品ID, 来自奖品中心
	PrizeName      string
	AllocatedStock int // 总预算库存 [cite: 141]
	Probability    float64
	IsSpecial      bool // 是否特殊奖品(如谢谢参与)
}
