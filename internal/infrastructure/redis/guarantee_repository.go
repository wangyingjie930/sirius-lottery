package redis

import (
	"context"
	"fmt"
	"github.com/wangyingjie930/nexus-pkg/redis"
	"sirius-lottery/internal/domain/strategy"
	"strconv"
)

// redisGuaranteeRepository 是 GuaranteeRepository 的 Redis 实现
type redisGuaranteeRepository struct {
	redisClient *redis.Client
	keyGen      *KeyGenerator
}

// NewRedisGuaranteeRepository 创建一个新的 GuaranteeRepository Redis 实现实例
func NewRedisGuaranteeRepository(redisClient *redis.Client) strategy.GuaranteeRepository {
	return &redisGuaranteeRepository{
		redisClient: redisClient,
		keyGen:      &KeyGenerator{},
	}
}

// IncrementAndGet 使用 Redis 的 INCR 命令原子性地增加用户的连续未中奖次数，并返回自增后的值。
// 如果 key 不存在，INCR 会先将其初始化为 0 再执行自增，所以第一次调用会返回 1。
func (r *redisGuaranteeRepository) IncrementAndGet(ctx context.Context, instanceID string, userID int64) (int, error) {
	// 从 key generator 获取标准化的 key
	key := r.keyGen.GuaranteeCounter(instanceID, strconv.FormatInt(userID, 10))

	// 执行 INCR 命令
	result, err := r.redisClient.GetClient().Incr(ctx, key).Result()
	if err != nil {
		// 如果 Redis 命令执行失败，返回错误
		return 0, fmt.Errorf("redis INCR for guarantee counter failed: %w", err)
	}

	return int(result), nil
}

// ResetCounter 当用户中奖后，使用 Redis 的 DEL 命令删除对应的保底计数器 key。
func (r *redisGuaranteeRepository) ResetCounter(ctx context.Context, instanceID string, userID int64) error {
	// 从 key generator 获取标准化的 key
	key := r.keyGen.GuaranteeCounter(instanceID, strconv.FormatInt(userID, 10))

	// 执行 DEL 命令
	if err := r.redisClient.GetClient().Del(ctx, key).Err(); err != nil {
		// 在生产环境中，这里的错误应该被记录到日志中。
		// 因为这只是一个重置操作，即使失败了，最坏的情况也只是用户下次抽奖时保底计数器没被清零。
		// 这比因为这个错误导致整个抽奖流程失败要好。所以我们只记录错误，不向上传递阻塞性错误。
		fmt.Printf("CRITICAL: Failed to reset guarantee counter for user %d in instance %s. Error: %v\n", userID, instanceID, err)
	}

	return nil
}
