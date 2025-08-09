package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql" // 导入mysql驱动
	"github.com/wangyingjie930/nexus-pkg/bootstrap"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"os"
	"sirius-lottery/internal/interfaces"
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
			if err := Init(context.Background()); err != nil {
				panic(err)
			}

			handler := interfaces.NewHttpHandler(lotterySrv)
			handler.RegisterRoutes(appCtx.Mux)
			logger.Logger.Printf("✅ Promotion service routes registered.")
		},
	})
}
