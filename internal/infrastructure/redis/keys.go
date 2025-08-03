package redis

import "fmt"

// Redis Key 命名规范文件

const (
	// ===== 活动配置缓存 =====
	// lottery:instance:{instance_id}
	// 存储序列化后的完整活动配置JSON
	KeyActivityConfig = "lottery:instance:%s"

	// ===== 分片库存 =====
	// lottery:stock:{instance_id}:{prize_id}:{shard_id}
	// 存储剩余库存数，使用DECR原子扣减
	KeyStock = "lottery:stock:%s:%s:%d"

	// ===== 分布式锁 =====
	// lock:draw:{user_id}:{instance_id}
	// 存储request_id，SET NX PX实现，Lua脚本安全释放
	KeyDrawLock = "lock:draw:%s:%s"

	// ===== 保底计数器 =====
	// lottery:guarantee:{instance_id}:{user_id}
	// 存储连续未中奖次数，使用INCR和DEL
	KeyGuaranteeCounter = "lottery:guarantee:%s:%s"

	// ===== API限流 =====
	// rate_limit:{api_path}:{user_id}
	// 存储当前请求数，令牌桶或固定窗口算法
	KeyRateLimit = "rate_limit:%s:%s"

	// ===== 用户中奖缓存 =====
	// lottery:win_log:{instance_id}:{user_id}
	// ZSet结构，score: timestamp, value: prize_id
	// 只用于展示，减轻DB查询压力
	KeyWinLog = "lottery:win_log:%s:%s"
)

// Redis Key 生成器函数
type KeyGenerator struct{}

func (kg *KeyGenerator) ActivityConfig(instanceID string) string {
	return fmt.Sprintf(KeyActivityConfig, instanceID)
}

func (kg *KeyGenerator) Stock(instanceID, prizeID string, shardID int) string {
	return fmt.Sprintf(KeyStock, instanceID, prizeID, shardID)
}

func (kg *KeyGenerator) DrawLock(userID, instanceID string) string {
	return fmt.Sprintf(KeyDrawLock, userID, instanceID)
}

func (kg *KeyGenerator) GuaranteeCounter(instanceID, userID string) string {
	return fmt.Sprintf(KeyGuaranteeCounter, instanceID, userID)
}

func (kg *KeyGenerator) RateLimit(apiPath, userID string) string {
	return fmt.Sprintf(KeyRateLimit, apiPath, userID)
}

func (kg *KeyGenerator) WinLog(instanceID, userID string) string {
	return fmt.Sprintf(KeyWinLog, instanceID, userID)
}
