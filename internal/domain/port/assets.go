package port

import "context"

type AssetsService interface {
	TryDeduct(ctx context.Context, userId int64, requestId string, cost int) error
	CancelDeduct(ctx context.Context, userId int64, requestId string) error
}
