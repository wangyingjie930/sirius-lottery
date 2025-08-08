package interfaces

import (
	"context"
	"fmt"
	"sirius-lottery/internal/application"
	"sirius-lottery/internal/pkg/eventbus"
)

type LotteryEventHandler struct {
	lotteryService application.LotteryService
}

func NewLotteryEventHandler(lotteryService application.LotteryService) *LotteryEventHandler {
	return &LotteryEventHandler{lotteryService: lotteryService}
}

func (h *LotteryEventHandler) HandleLotteryWinEvent(ctx context.Context, event eventbus.Event) error {
	winEvent, ok := event.(*application.LotteryWinEvent)
	if !ok {
		return nil
	}

	// 1. Deduct Stock
	stockReq := &application.StockActionRequest{
		InstanceID: winEvent.WinRecord.InstanceID,
		PrizeID:    winEvent.WinRecord.PrizeID,
		Num:        1,
	}
	if err := h.lotteryService.DeductStock(ctx, stockReq); err != nil {
		// In a real system, you might want to retry or send a notification
		fmt.Printf("Failed to deduct stock: %v\n", err)
		// Optionally, you can try to compensate by increasing the stock back
		h.lotteryService.IncreaseStock(ctx, stockReq)
		return err
	}

	// 2. Deduct Assets (e.g., points, coins)
	// In a real microservices architecture, this would be a call to an asset service.
	// Here, we'll just log it.
	assetReq := &application.AssetRequest{
		UserId: int64(winEvent.WinRecord.UserID),
		// You might need to get cost information from the win record or another source
		// Cost: winEvent.WinRecord.Cost,
	}
	fmt.Printf("Deducting assets for user %d, request: %+v\n", assetReq.UserId, assetReq)

	// 3. Update win record status
	// This would typically be done by updating the record in the database
	fmt.Printf("Updating win record %s to SUCCESS\n", winEvent.WinRecord.OrderID)

	return nil
}

func (h *LotteryEventHandler) Register(bus *eventbus.MemoryEventBus) {
	bus.Subscribe(application.TopicLotteryWin, h.HandleLotteryWinEvent)
}
