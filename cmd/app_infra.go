package main

import (
	"context"
	"github.com/wangyingjie930/nexus-pkg/bootstrap"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"github.com/wangyingjie930/nexus-pkg/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sirius-lottery/internal/infrastructure/eventbus"
	model "sirius-lottery/internal/infrastructure/gorm"
)

type AppDependencies struct {
	DB          *gorm.DB
	AppEventbus eventbus.Producer
	Redis       *redis.Client
}

func InitDependencies(ctx context.Context) (*AppDependencies, error) {
	appDependencies := &AppDependencies{
		DB:          InitDb(),
		Redis:       InitRedis(),
		AppEventbus: InitEventbus(),
	}
	return appDependencies, nil
}

func InitDb() *gorm.DB {
	config := bootstrap.GetCurrentConfig()
	dsn := config.Infra.Mysql.Addrs // 应从配置获取
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Error().Err(err).Msgf("failed to connect to database with gorm: %v", err)
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 2. **自动迁移 (基础设施)**
	// 使用在 infrastructure 包中定义的 GORM 模型
	err = db.AutoMigrate(
		&model.LotteryTemplate{},
		&model.LotteryInstance{},
		&model.LotteryPool{},
		&model.LotteryWinRecord{}, // Depends on LotteryInstance
		&model.LotteryPrize{},     // Depends on LotteryPool
	)
	if err != nil {
		logger.Logger.Error().Err(err).Msgf("WARN: failed to auto migrate gorm models: %v", err)
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 1") // 3. **创建仓储实例 (基础设施)**
	return db
}

func InitEventbus() eventbus.Producer {
	config := bootstrap.GetCurrentConfig()

	producer, err := eventbus.NewProducer(config.Infra.Kafka.Brokers, "sirius-lottery", "", 1)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("InitEventbus")
	}

	return producer
}

func InitRedis() *redis.Client {
	redisCilent, err := redis.NewClient(bootstrap.GetCurrentConfig().Infra.Redis.Addrs)
	if err != nil {
		panic(err)
	}
	return redisCilent
}
