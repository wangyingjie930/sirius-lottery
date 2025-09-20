package main

import (
	"context"
	"github.com/wangyingjie930/nexus-pkg/bootstrap"
	"sirius-lottery/internal/application"
	strategy2 "sirius-lottery/internal/domain/strategy"
	"sirius-lottery/internal/infrastructure"
	"sirius-lottery/internal/infrastructure/eventbus"
	"sirius-lottery/internal/infrastructure/port"
	redis2 "sirius-lottery/internal/infrastructure/redis"
	"sirius-lottery/internal/infrastructure/repository"
)

var (
	lotterySrv application.LotteryService
)

func Init(ctx context.Context) (err error) {
	dept, err := InitDependencies(ctx)
	if err != nil {
		return err
	}

	lrp := repository.NewGormLotteryRepository(dept.DB, dept.Redis)
	wrp := repository.NewGormWinRecordRepository(dept.DB)
	gp := redis2.NewRedisGuaranteeRepository(dept.Redis)
	strategy := strategy2.NewLotteryStrategyFactory(gp)
	uow := infrastructure.NewGormUnitOfWork(dept.DB)

	assetSrv := port.NewAssetSrv()
	stockSrv := port.NewStockSrv()

	l := application.NewLotteryServiceImpl(lrp, wrp, strategy, uow, dept.AppEventbus, assetSrv, stockSrv)
	config := bootstrap.GetCurrentConfig()
	if err = eventbus.NewConsumerService().RegisterConsumer(config.Infra.Kafka.Brokers, "sirius-lottery", "sirius-lottery", l); err != nil {
		return err
	}

	lotterySrv = l

	return nil
}
