package main

import (
	_ "github.com/go-sql-driver/mysql" // 导入mysql驱动
	"github.com/wangyingjie930/nexus-pkg/bootstrap"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"github.com/wangyingjie930/nexus-pkg/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sirius-lottery/internal/application"
	strategy2 "sirius-lottery/internal/domain/strategy"
	"sirius-lottery/internal/infrastructure"
	model "sirius-lottery/internal/infrastructure/gorm"
	redis2 "sirius-lottery/internal/infrastructure/redis"
	"sirius-lottery/internal/infrastructure/repository"
	"sirius-lottery/internal/interfaces"
	"sirius-lottery/internal/pkg/eventbus"
)

const (
	serviceName = "lottery-service"
)

func main() {
	os.Setenv("NEXUS_CONFIG_PATH", "./config/local.yaml")

	// 初始化 tracer, nacos 等通用组件
	bootstrap.Init()

	bootstrap.StartService(bootstrap.AppInfo{
		ServiceName: serviceName,
		Port:        8080,
		RegisterHandlers: func(appCtx bootstrap.AppCtx) {
			// 1. **连接数据库 (基础设施)**
			// dsn := bootstrap.GetCurrentConfig().DB.Source
			dsn := bootstrap.GetCurrentConfig().Infra.Mysql.Addrs // 应从配置获取
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
			//couponRepository := infrastructure.NewGormCouponRepository(db)
			//templateRepo := infrastructure.NewGormPromotionTemplateRepository(db)
			redisCilent, _ := redis.NewClient(bootstrap.GetCurrentConfig().Infra.Redis.Addrs)
			lrp := repository.NewGormLotteryRepository(db, redisCilent)
			wrp := repository.NewGormWinRecordRepository(db)
			gp := redis2.NewRedisGuaranteeRepository(redisCilent)
			strategy := strategy2.NewLotteryStrategyFactory(gp)
			uow := infrastructure.NewGormUnitOfWork(db)

			// 4. **创建应用服务实例 (应用层)**
			// 将仓储接口注入到应用服务中
			//tracer := otel.Tracer(serviceName)
			eventBus := eventbus.NewMemoryEventBus()
			lotterySrv := application.NewLotteryServiceImpl(lrp, wrp, strategy, uow, eventBus)

			// 5. **创建HTTP处理器 (接口层)**
			// 将应用服务注入到HTTP处理器中
			httpHandler := interfaces.NewHttpHandler(lotterySrv)

			// 6. **创建事件处理器并注册**
			eventHandler := interfaces.NewLotteryEventHandler(lotterySrv)
			eventHandler.Register(eventBus)

			// 7. **启动服务并注册路由**
			httpHandler.RegisterRoutes(appCtx.Mux)

			logger.Logger.Printf("✅ Promotion service routes registered.")
		},
	})
}
