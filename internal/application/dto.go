package application

import "time"

// DrawRequest 对应 /draw 接口的请求体 [cite: 156]
type DrawRequest struct {
	InstanceID string `json:"instance_id"`
	RequestID  string `json:"request_id"`
}

// DrawResponse 对应 /draw 接口的成功响应 [cite: 157]
type DrawResponse struct {
	OrderID string `json:"order_id"`
	PrizeID string `json:"prize_id"`
	IsWin   bool   `json:"is_win"`
}

type LotteryInstanceResponse struct {
	InstanceId     string    `json:"instance_id"`
	Name           string    `json:"name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	ServerTime     time.Time `json:"server_time"`
	TemplateStyle  string    `json:"template_style"`
	TemplateConfig struct {
		BackgroundImage string `json:"background_image"`
	} `json:"template_config"`
	Pools []struct {
		PoolId   string `json:"pool_id"`
		PoolName string `json:"pool_name"`
		Cost     []struct {
			AssetId string `json:"asset_id"`
			Amount  int    `json:"amount"`
		} `json:"cost"`
		Prizes []struct {
			Position int    `json:"position"`
			PrizeId  string `json:"prize_id"`
		} `json:"prizes"`
	} `json:"pools"`
}

type StockActionRequest struct {
	InstanceID string `json:"instance_id"`
	PrizeID    string `json:"prize_id"`
	Num        int    `json:"num"`
}

type CreateWinRecordRequest struct {
	OrderID    string `json:"order_id"`
	PrizeID    string `json:"prize_id"`
	InstanceID string `json:"instance_id"`
	RequestID  string `json:"request_id"`
	UserID     int64  `json:"user_id"`
}

type AssetRequest struct {
	Cost   int   `json:"cost"`
	UserId int64 `json:"user_id"`
}
