package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/wangyingjie930/nexus-pkg/transactional"
	"gorm.io/gorm"
	"sirius-lottery/internal/domain"
	"sirius-lottery/internal/domain/entity"
	"sirius-lottery/internal/domain/port"
	"sirius-lottery/internal/domain/strategy"

	"github.com/dtm-labs/client/dtmcli"
	"time"
)

type lotteryServiceImpl struct {
	repo             domain.LotteryRepository
	winRecordRepo    domain.WinRecordRepository
	locker           domain.Locker
	assetsSrv        port.AssetsService
	strategyFact     *strategy.LotteryStrategyFactory
	transactionalSrv *transactional.Service
	uow              domain.UnitOfWork
}

const (
	DtmServer        = "http://localhost:36789/api/dtmsvr"
	AssetsServiceURL = ""
)

// Draw 实现了核心抽奖逻辑
func (s *lotteryServiceImpl) Draw(ctx context.Context, req *DrawRequest) (*DrawResponse, error) {
	// 从 context 中获取 userID, 这里假设 userID 已经通过上游中间件注入
	// 在实际项目中, 通常会使用 JWT 或者其他 session 机制来获取用户信息
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		return nil, errors.New("无法获取用户信息")
	}

	// === 步骤 1: 获取活动配置缓存 (Redis) ===
	// repo.GetInstance 应该优先从缓存读取
	instance, err := s.repo.GetInstance(ctx, req.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("获取活动配置失败: %w", err)
	}

	// === 步骤 2: 前置校验 (活动状态/时间/人群) ===
	if err := instance.Check(time.Now()); err != nil {
		return nil, err
	}
	// 人群校验 (UserScope)
	if !instance.IsUserAllowed(userID) {
		return nil, errors.New("您不符合参与条件")
	}

	// === 步骤 3: 获取分布式锁
	//if !s.locker.Lock(userID, instance.InstanceID) {
	//	return nil, errors.New("操作频繁，请稍后再试")
	//}
	//defer s.locker.UnLock(userID, instance.InstanceID) // 确保锁最终被释放

	// === 步骤 4: [幂等性检查] 查询中奖/消息记录表，检查 req_id 是否已处理 ===
	// 这里我们假设如果一个请求ID已经存在于中奖记录中，就认为是重复请求
	//existingRecord, err := s.winRecordRepo.GetByRequestID(ctx, req.RequestID)
	//if err != nil {
	//	// 查询出错，需要假定幂等性检查失败，但不能认为是重复请求
	//	return nil, fmt.Errorf("幂等性检查失败: %w", err)
	//}
	//if existingRecord != nil {
	//	// 找到记录，说明是重复请求，返回上次的结果
	//	return &DrawResponse{
	//		OrderID: existingRecord.OrderID,
	//		PrizeID: existingRecord.PrizeID,
	//		IsWin:   !existingRecord.IsThankYouPrize(), // 假设我们有一个方法判断是否为谢谢参与
	//	}, nil
	//}

	// === 步骤 5: [RPC] 扣减资产 (含幂等ID) ===
	// 假设抽奖消耗定义在第一个奖池中
	// TODO: 在多奖池场景下，需要明确指定从哪个奖池抽奖
	if len(instance.Pools) == 0 {
		return nil, errors.New("活动配置错误：缺少奖池")
	}
	cost := instance.Pools[0].GetCost() // 假设 GetCost 返回需要扣减的资产信息

	gid := dtmcli.MustGenGid(DtmServer)
	var drawResp *DrawResponse
	err := dtmcli.TccGlobalTransaction(DtmServer, gid, func(tcc *dtmcli.Tcc) (*resty.Response, error) {
		resp, err := tcc.CallBranch(nil, AssetsServiceURL+"/TryDeduct", AssetsServiceURL+"/DeductConfirm", AssetsServiceURL+"/DeductRevert")
		if err != nil {
			return resp, err
		}

		pool := instance.Pools[0]
		strategy, err := s.strategyFact.GetStrategy(pool.LotteryStrategy)
		if err != nil {
			return nil, fmt.Errorf("抽奖策略加载失败: %w", err)
		}

		drawCtx := &domain.DrawContext{
			InstanceID: instance.InstanceID,
			UserID:     userID,
			Pool:       &pool,
			Prizes:     pool.Prizes,
		}
		wonPrize, err := strategy.Draw(ctx, drawCtx)
		if err != nil {
			return nil, fmt.Errorf("抽奖执行失败: %w", err)
		}

		if wonPrize.IsSpecial {
			drawResp = &DrawResponse{
				OrderID: "THANK_YOU_ORDER", // 可以给一个特殊订单号
				PrizeID: wonPrize.PrizeID,
				IsWin:   false,
			}
			return nil, nil
		}

		drawResp = &DrawResponse{
			OrderID: uuid.New().String(),
			PrizeID: wonPrize.PrizeID,
			IsWin:   !wonPrize.IsSpecial,
		}

		resp, err = tcc.CallBranch(nil, "/TryDeductStock", "/ConfirmDeductStock", "/cancelDeductStock")
		if err != nil {
			return resp, err
		}

		return tcc.CallBranch(nil, "/TryRecord", "/ConfirmRecord", "/cancelRecord")
	})

	if err != nil {
		return nil, err
	}

	return drawResp, nil
	//if err := s.assetsSrv.TryDeduct(ctx, userID, req.RequestID, cost); err != nil {
	//	// 资产扣减失败，无需归还，因为 TryDeduct 应该已经是幂等的
	//	return nil, fmt.Errorf("资产不足或扣减失败: %w", err)
	//}

	// === 步骤 6: 执行抽奖算法 (内存计算, 得到目标prize_id) ===
	// 同样，暂时只考虑第一个奖池

	// === 步骤 7: [Lua] 扣减分片库存 ===
	// 如果中的不是"谢谢参与"等特殊奖品，才需要扣库存
	//if !wonPrize.IsSpecial {
	//	stockDeducted, err := s.repo.DeductStock(ctx, instance.InstanceID, wonPrize.PrizeID, 1)
	//	if err != nil || !stockDeducted {
	//		// 库存不足或扣减失败，归还资产
	//		s.assetsSrv.CancelDeduct(ctx, userID, req.RequestID) // 补偿操作
	//		return nil, errors.New("非常遗憾，奖品已被抢光了")
	//	}
	//}
	//
	//// === 步骤 8: [DB事务] 写中奖记录 + 写发奖消息(Outbox) ===
	//orderID := uuid.New().String() // 生成唯一订单ID
	//winRecord := &entity.LotteryWinRecord{
	//	OrderID:    orderID,
	//	RequestID:  req.RequestID, // 存储请求ID用于幂等性判断
	//	InstanceID: instance.InstanceID,
	//	UserID:     uint64(userID),
	//	PrizeID:    wonPrize.PrizeID,
	//	Status:     entity.WinRecordStatusPending, // 初始状态为待发放
	//}
	//
	//msgPayload, _ := json.Marshal(winRecord)
	//// 创建事务性消息 (Outbox Pattern)
	//
	//// 使用事务管理器执行数据库操作
	//err = s.uow.Execute(ctx, func(rp domain.RepositoryProvider) error {
	//	if err := rp.WinRecordRepository().Create(ctx, winRecord); err != nil {
	//		return err
	//	}
	//	if err := rp.TransactionalStore().CreateInTx(ctx, &transactional.Message{
	//		Topic:   "lottery_win_events",
	//		Key:     uuid.New().String(),
	//		Payload: msgPayload,
	//		Status:  transactional.StatusPending,
	//	}); err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//
	//if err != nil {
	//	// DB事务失败，这是最关键的失败场景
	//	// 1. 补偿库存 (非常重要)
	//	if !wonPrize.IsSpecial {
	//		s.repo.IncreaseStock(ctx, instance.InstanceID, wonPrize.PrizeID, 1)
	//	}
	//	// 2. 补偿资产
	//	s.assetsSrv.CancelDeduct(ctx, userID, req.RequestID)
	//	return nil, fmt.Errorf("系统繁忙，请稍后再试: %w", err)
	//}
	//
	//// === 步骤 9: 释放锁 ===
	//// 由 defer s.locker.UnLock(lockKey) 自动完成
	//
	//// === 步骤 10: [异步] 消息通过CDC投递 ===
	//// 此步骤由数据库和CDC工具链保证，应用代码层面无需处理
	//
	//// === 返回中奖结果 ===
	//return &DrawResponse{
	//	OrderID: orderID,
	//	PrizeID: wonPrize.PrizeID,
	//	IsWin:   !wonPrize.IsSpecial,
	//}, nil
}

func (s *lotteryServiceImpl) DeductStock(ctx context.Context, req *StockActionRequest) error {
	deducted, err := s.repo.DeductStock(ctx, req.InstanceID, req.PrizeID, req.Num)
	if err != nil {
		return err // 返回错误，DTM会重试
	}
	if !deducted {
		// 库存不足是业务失败，应返回 dtmcli.ErrFailure，
		// 这会立即中止SAGA并触发回滚，不会重试。
		return dtmcli.ErrFailure
	}
	return nil
}

func (s *lotteryServiceImpl) IncreaseStock(ctx context.Context, req *StockActionRequest) error {
	// 补偿操作必须保证成功。如果底层调用(如Redis)失败，应记录严重错误日志并告警，
	// 以便人工介入，但需要向 DTM 返回成功，以允许其他分支的补偿继续执行。
	if _, err := s.repo.IncreaseStock(ctx, req.InstanceID, req.PrizeID, req.Num); err != nil {
		fmt.Printf("CRITICAL: SAGA compensation 'IncreaseStock' failed. Request: %+v. Error: %v\n", req, err)
	}
	return nil
}

func (s *lotteryServiceImpl) CreateWinRecordInTx(ctx context.Context, req *CreateWinRecordRequest, tx *gorm.DB) error {
	// 此方法现在接收一个 *gorm.DB 对象，这个事务由 DTM Barrier 管理。
	// 我们需要一个使用这个特定事务的 RepositoryProvider。
	repoProvider := infrastructure.NewGormRepoProvider(tx)

	// 1. 创建中奖记录
	winRecord := &entity.LotteryWinRecord{
		OrderID:    req.OrderID,
		RequestID:  req.RequestID,
		InstanceID: req.InstanceID,
		UserID:     uint64(req.UserID),
		PrizeID:    req.PrizeID,
		Status:     entity.WinRecordStatusPending,
	}
	if err := repoProvider.WinRecordRepository().Create(ctx, winRecord); err != nil {
		return fmt.Errorf("创建中奖记录失败: %w", err)
	}

	// 2. 使用 Outbox 模式创建发奖消息
	msgPayload, _ := json.Marshal(winRecord)
	msg := &transactional.Message{
		Topic:   "lottery_win_events",
		Key:     uuid.New().String(),
		Payload: msgPayload,
		Status:  transactional.StatusPending,
	}
	if err := repoProvider.TransactionalStore().CreateInTx(ctx, msg); err != nil {
		return fmt.Errorf("创建发奖消息失败: %w", err)
	}

	return nil
}
