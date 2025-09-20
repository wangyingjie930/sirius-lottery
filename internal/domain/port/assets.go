package port

import "context"

type StockActionRequest struct {
	InstanceID string `json:"instance_id"`
	PrizeID    string `json:"prize_id"`
	Num        int    `json:"num"`
}

type AssetsService interface {
	ActionName() string
	ComponentName() string
	TryDeduct(ctx context.Context, request StockActionRequest) error
	CancelDeduct(ctx context.Context, request StockActionRequest) error
	ConfirmDeduct(ctx context.Context, request StockActionRequest) error
}

type StockService interface {
	ActionName() string
	ComponentName() string
	TryDeduct(ctx context.Context, request StockActionRequest) error
	CancelDeduct(ctx context.Context, request StockActionRequest) error
	ConfirmDeduct(ctx context.Context, request StockActionRequest) error
}
