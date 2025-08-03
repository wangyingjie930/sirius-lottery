package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	redisClient "github.com/wangyingjie930/nexus-pkg/redis"
	"gorm.io/gorm"
	"sirius-lottery/internal/domain/entity"
	gorm_model "sirius-lottery/internal/infrastructure/gorm"
	redis_key "sirius-lottery/internal/infrastructure/redis"
	"time"
)

type gormLotteryRepository struct {
	db     *gorm.DB
	redis  *redisClient.Client
	keyGen *redis_key.KeyGenerator
}

// NewGormLotteryRepository 创建一个新的 LotteryRepository GORM 实现
func NewGormLotteryRepository(db *gorm.DB, redis *redisClient.Client) *gormLotteryRepository {
	return &gormLotteryRepository{
		db:     db,
		redis:  redis,
		keyGen: &redis_key.KeyGenerator{},
	}
}

// GetInstance 优先从 Redis 缓存获取活动实例，如果缓存未命中，则从数据库加载并回写缓存。
func (r *gormLotteryRepository) GetInstance(ctx context.Context, instanceId string) (*entity.LotteryInstance, error) {
	// 1. 从 Redis 获取缓存
	key := r.keyGen.ActivityConfig(instanceId)
	val, err := r.redis.GetClient().Get(ctx, key).Result()

	if err == nil {
		// 缓存命中
		var instance entity.LotteryInstance
		if err := json.Unmarshal([]byte(val), &instance); err == nil {
			return &instance, nil
		}
		// JSON 解码失败，记录日志，继续往下走从数据库加载
		fmt.Printf("Error unmarshalling instance from redis: %v\n", err)
	} else if err != redis.Nil {
		// Redis 查询出错
		return nil, fmt.Errorf("failed to get instance from redis: %w", err)
	}

	// 2. 缓存未命中，从数据库加载
	var instanceModel gorm_model.LotteryInstance
	err = r.db.WithContext(ctx).
		Preload("Pools.Prizes"). // 预加载奖池和奖品
		Where("instance_id = ?", instanceId).
		First(&instanceModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lottery instance not found")
		}
		return nil, fmt.Errorf("failed to get instance from db: %w", err)
	}

	// 3. 将 GORM Model 转换为领域实体 (Domain Entity)
	domainInstance := toDomainInstance(&instanceModel)

	// 4. 回写 Redis 缓存
	instanceBytes, err := json.Marshal(domainInstance)
	if err == nil {
		// 设置缓存，并给一个合理的过期时间，比如活动结束后的一段时间
		ttl := time.Until(domainInstance.EndTime.Add(1 * time.Hour))
		if ttl < 0 {
			ttl = 5 * time.Minute // 如果活动已结束，给一个短的缓存时间
		}
		r.redis.GetClient().Set(ctx, key, instanceBytes, ttl).Err()
	}

	ttl := time.Until(domainInstance.EndTime.Add(1 * time.Hour))
	for _, pools := range domainInstance.Pools {
		for _, prize := range pools.Prizes {
			key := r.keyGen.Stock(instanceId, prize.PrizeID, 0)
			r.redis.GetClient().Set(ctx, key, prize.AllocatedStock, ttl)
		}
	}

	return domainInstance, nil
}

// CheckIdempotencyKey 检查幂等键是否存在 (例如，防止重复抽奖)
func (r *gormLotteryRepository) CheckIdempotencyKey(ctx context.Context, key string) bool {
	// SET NX a a PX 300000
	// 使用 Redis 的 SET NX (Set if Not Exists) 来原子性地检查和设置键
	// 如果键不存在，设置成功并返回 true (表示是第一次请求)
	// 如果键已存在，设置失败并返回 false (表示是重复请求)
	// 给键设置一个过期时间，以防应用崩溃导致锁无法释放
	ok, err := r.redis.GetClient().SetNX(ctx, key, "1", 5*time.Minute).Result()
	if err != nil {
		// 在生产环境中，这里应该有更健壮的错误处理和日志记录
		fmt.Printf("Error checking idempotency key in redis: %v\n", err)
		return false // 出错时默认为重复请求，防止潜在的重复处理
	}
	return ok
}

// DeductStock 使用 Redis DECR 原子性地扣减库存
func (r *gormLotteryRepository) DeductStock(ctx context.Context, instanceId string, prizeId string, num int) (bool, error) {
	// TODO: 目前是单点库存，实际大厂会使用分片库存来提高性能
	// 使用 Lua 脚本保证操作的原子性
	key := r.keyGen.Stock(instanceId, prizeId, 0) // shard_id 暂时为0

	// Lua 脚本:
	// 1. 获取当前库存
	// 2. 如果库存足够，则扣减
	// 3. 返回扣减后的结果 (1 表示成功, 0 表示库存不足)
	script := `
        local stock = redis.call('GET', KEYS[1])
        if not stock or tonumber(stock) < tonumber(ARGV[1]) then
            return 0
        end
        redis.call('DECRBY', KEYS[1], ARGV[1])
        return 1
    `
	res, err := r.redis.GetClient().Eval(ctx, script, []string{key}, num).Result()
	if err != nil {
		return false, fmt.Errorf("failed to execute deduct stock lua script: %w", err)
	}

	if res.(int64) == 0 {
		return false, nil // 库存不足
	}

	return true, nil // 扣减成功
}

// IncreaseStock 使用 Redis INCR 原子性地增加库存 (用于取消/回滚)
func (r *gormLotteryRepository) IncreaseStock(ctx context.Context, instanceId string, prizeId string, num int) (bool, error) {
	key := r.keyGen.Stock(instanceId, prizeId, 0) // shard_id 暂时为0
	err := r.redis.GetClient().IncrBy(ctx, key, int64(num)).Err()
	if err != nil {
		return false, fmt.Errorf("failed to increase stock in redis: %w", err)
	}
	return true, nil
}

// toDomainInstance 将 GORM 模型转换为领域实体
func toDomainInstance(m *gorm_model.LotteryInstance) *entity.LotteryInstance {
	instance := &entity.LotteryInstance{
		ID:           m.ID,
		InstanceID:   m.InstanceID,
		InstanceName: m.InstanceName,
		TemplateID:   m.TemplateID,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		Status:       m.Status,
		Pools:        make([]entity.LotteryPool, 0, len(m.Pools)),
	}
	for _, p := range m.Pools {
		pool := entity.LotteryPool{
			ID:              int64(p.ID),
			InstanceID:      p.InstanceID,
			PoolName:        p.PoolName,
			LotteryStrategy: p.LotteryStrategy,
			Prizes:          make([]*entity.LotteryPrize, 0, len(p.Prizes)),
		}
		// CostJSON 和 StrategyConfigJSON 的转换
		costBytes, _ := p.CostJSON.Value()
		if bs, ok := costBytes.([]byte); ok {
			pool.CostJSON = string(bs)
		}
		strategyBytes, _ := p.StrategyConfigJSON.Value()
		if bs, ok := strategyBytes.([]byte); ok {
			pool.StrategyConfigJSON = string(bs)
		}

		for _, pz := range p.Prizes {
			prize := &entity.LotteryPrize{
				ID:             int64(pz.ID),
				PoolID:         int64(pz.PoolID),
				PrizeID:        pz.PrizeID,
				PrizeName:      pz.PrizeName,
				AllocatedStock: pz.AllocatedStock, //
				Probability:    pz.Probability,
				IsSpecial:      pz.IsSpecial,
			}
			pool.Prizes = append(pool.Prizes, prize)
		}
		instance.Pools = append(instance.Pools, pool)
	}
	return instance
}
