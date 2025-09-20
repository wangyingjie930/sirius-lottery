package port

import (
	"context"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"sirius-lottery/internal/domain/port"
)

const (
	AssetsServiceURL  = "http://host.docker.internal:8080/api/v2/lottery/dtm"
	LotteryServiceURL = "http://host.docker.internal:8080/api/v2/lottery/dtm"
	LocalHost         = "http://host.docker.internal:8080/api/v2/lottery/dtm"
)

type AssetSrv struct {
}

func NewAssetSrv() *AssetSrv {
	return &AssetSrv{}
}

func (a *AssetSrv) ActionName() string {
	return AssetsServiceURL + "/asset/try"
}

func (a *AssetSrv) ComponentName() string {
	return AssetsServiceURL + "/asset/cancel"
}

func (a *AssetSrv) TryDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ AssetTry")
	return nil
}

func (a *AssetSrv) CancelDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ CancelDeduct")
	return nil
}

func (a *AssetSrv) ConfirmDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ ConfirmDeduct")
	return nil
}

type StockSrv struct {
}

func (s *StockSrv) ActionName() string {
	return LotteryServiceURL + "/stock/try"
}

func (s *StockSrv) ComponentName() string {
	return LotteryServiceURL + "/stock/cancel"
}

func (s *StockSrv) TryDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ TryDeduct")
	return nil
}

func (s *StockSrv) CancelDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ CancelDeduct")
	return nil
}

func (s *StockSrv) ConfirmDeduct(ctx context.Context, request port.StockActionRequest) error {
	logger.Ctx(ctx).Println("✅ ConfirmDeduct")
	return nil
}

func NewStockSrv() *StockSrv {
	return &StockSrv{}
}
