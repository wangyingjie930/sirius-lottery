package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wangyingjie930/nexus-pkg/transactional"
	"sirius-lottery/internal/domain"
	"sirius-lottery/internal/domain/entity"
	"sirius-lottery/internal/domain/strategy"
	"time"
)

type lotteryServiceImpl struct {
	repo          domain.LotteryRepository
	winRecordRepo domain.WinRecordRepository
	strategyFact  *strategy.LotteryStrategyFactory
	uow           domain.UnitOfWork
	eventBus      EventBus
}

func NewLotteryServiceImpl(repo domain.LotteryRepository, winRecordRepo domain.WinRecordRepository, strategyFact *strategy.LotteryStrategyFactory, uow domain.UnitOfWork, eventBus EventBus) *lotteryServiceImpl {
	return &lotteryServiceImpl{repo: repo, winRecordRepo: winRecordRepo, strategyFact: strategyFact, uow: uow, eventBus: eventBus}
}

// Draw 实现了核心抽奖逻辑
func (s *lotteryServiceImpl) Draw(ctx context.Context, req *DrawRequest) (*DrawResponse, error) {
	// 从 context 中获取 userID, 这里假设 userID 已经通过上游中间件注入
	// 在实际项目中, 通常会使用 JWT 或者其他 session 机制来获取用户信息
	var userID int64 = 100

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

	// TODO: 在多奖池场景下，需要明确指定从哪个奖池抽奖
	if len(instance.Pools) == 0 {
		return nil, errors.New("活动配置错误：缺少奖池")
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
		return &DrawResponse{
			OrderID: "THANK_YOU_ORDER", // 可以给一个特殊订单号
			PrizeID: wonPrize.PrizeID,
			IsWin:   false,
		}, nil
	}

	drawResp := &DrawResponse{
		OrderID: uuid.New().String(),
		PrizeID: wonPrize.PrizeID,
		IsWin:   !wonPrize.IsSpecial,
	}

	winRecord := &entity.LotteryWinRecord{
		OrderID:    drawResp.OrderID,
		RequestID:  req.RequestID,
		InstanceID: req.InstanceID,
		UserID:     uint64(userID),
		PrizeID:    drawResp.PrizeID,
		Status:     entity.WinRecordStatusPending,
	}

	err = s.uow.Execute(ctx, func(repoProvider domain.RepositoryProvider) error {
		// 1. 创建中奖记录
		if err := repoProvider.WinRecordRepository().Create(ctx, winRecord); err != nil {
			return fmt.Errorf("创建中奖记录失败: %w", err)
		}

		// 2. 使用 Outbox 模式创建发奖消息
		msgPayload, _ := json.Marshal(winRecord)
		msg := &transactional.Message{
			Topic:   TopicLotteryWin,
			Key:     uuid.New().String(),
			Payload: msgPayload,
			Status:  transactional.StatusPending,
		}
		if err := repoProvider.TransactionalStore().CreateInTx(ctx, msg); err != nil {
			return fmt.Errorf("创建发奖消息失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish event
	s.eventBus.Publish(ctx, &LotteryWinEvent{WinRecord: winRecord})

	return drawResp, nil
}

func (s *lotteryServiceImpl) DeductStock(ctx context.Context, req *StockActionRequest) error {
	deducted, err := s.repo.DeductStock(ctx, req.InstanceID, req.PrizeID, req.Num)
	if err != nil {
		return err
	}
	if !deducted {
		return errors.New("stock not enough")
	}
	return nil
}

func (s *lotteryServiceImpl) IncreaseStock(ctx context.Context, req *StockActionRequest) error {
	if _, err := s.repo.IncreaseStock(ctx, req.InstanceID, req.PrizeID, req.Num); err != nil {
		fmt.Printf("CRITICAL: compensation 'IncreaseStock' failed. Request: %+v. Error: %v\n", req, err)
	}
	return nil
}

func (s *lotteryServiceImpl) GetLotteryInstance(ctx context.Context, instanceID string) (*LotteryInstanceResponse, error) {
	return nil, nil
}
