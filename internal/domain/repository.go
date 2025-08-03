package domain

import (
	"context"
	"github.com/wangyingjie930/nexus-pkg/transactional"
	"sirius-lottery/internal/domain/entity"
)

type LotteryRepository interface {
	GetInstance(ctx context.Context, instanceId string) (*entity.LotteryInstance, error)
	CheckIdempotencyKey(ctx context.Context, key string) bool
	DeductStock(ctx context.Context, instanceId string, prizeId string, num int) (bool, error)
	IncreaseStock(ctx context.Context, instanceId string, prizeId string, num int) (bool, error)
}

//type Locker interface {
//	Lock(userID int64, instanceID string) bool
//	UnLock(userID int64, instanceID string)
//}

// WinRecordRepository 定义中奖记录的数据库操作
type WinRecordRepository interface {
	Create(ctx context.Context, record *entity.LotteryWinRecord) error
	GetByRequestID(ctx context.Context, requestID string) (*entity.LotteryWinRecord, error)
}

// UnitOfWork 定义了工作单元的接口
// 它提供了一种方式来确保多个仓储操作在同一个事务中执行
type UnitOfWork interface {
	// Execute 将一个函数包裹在单个事务中执行
	// fn 是包含所有业务逻辑和仓储操作的函数
	// 如果 fn 返回错误，事务将回滚；否则，事务将提交
	Execute(ctx context.Context, fn func(repoProvider RepositoryProvider) error) error
}

// RepositoryProvider 是一个接口，用于在事务中获取仓储实例
// 这样可以确保所有获取到的仓储都共享同一个事务
type RepositoryProvider interface {
	LotteryRepository() LotteryRepository
	WinRecordRepository() WinRecordRepository
	TransactionalStore() transactional.Store
}
