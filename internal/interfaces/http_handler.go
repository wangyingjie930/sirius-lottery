package interfaces

import (
	"context"
	"encoding/json"
	"github.com/dtm-labs/client/dtmcli"
	"net/http"
	"sirius-lottery/internal/application"
	port2 "sirius-lottery/internal/domain/port"
	"sirius-lottery/internal/infrastructure/port"
)

type HttpHandler struct {
	lotteryService application.LotteryService

	assetSrv port.AssetSrv
	stockSrv port.StockSrv
}

func NewHttpHandler(lotteryService application.LotteryService) *HttpHandler {
	return &HttpHandler{lotteryService: lotteryService}
}

func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v2/lottery/draw", h.Draw)                               //
	mux.HandleFunc("GET /api/v2/lottery/instance/{instanceId}", h.GetLotteryInstance) //

	// DTM TCC Handlers
	mux.HandleFunc("POST /api/v2/lottery/dtm/stock/try", h.StockDeductTry)
	mux.HandleFunc("POST /api/v2/lottery/dtm/stock/confirm", h.StockDeductConfirm)
	mux.HandleFunc("POST /api/v2/lottery/dtm/stock/cancel", h.StockDeductCancel)

	mux.HandleFunc("POST /api/v2/lottery/dtm/asset/try", h.AssetTry)
	mux.HandleFunc("POST /api/v2/lottery/dtm/asset/confirm", h.AssetConfirm)
	mux.HandleFunc("POST /api/v2/lottery/dtm/asset/cancel", h.AssetCancel)
}

func (h *HttpHandler) Draw(writer http.ResponseWriter, request *http.Request) {
	var req application.DrawRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// 在实际项目中，UserID 应该从认证中间件（如JWT）中获取
	ctx := context.WithValue(request.Context(), "userID", int64(100)) // 示例 userID

	resp, err := h.lotteryService.Draw(ctx, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}

func (h *HttpHandler) GetLotteryInstance(writer http.ResponseWriter, request *http.Request) {
	instanceId := request.PathValue("instanceId")
	if instanceId == "" {
		http.Error(writer, "instanceId is required", http.StatusBadRequest)
		return
	}

	resp, err := h.lotteryService.GetLotteryInstance(request.Context(), instanceId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}

// writeDtmResponse 封装了向 DTM 返回结果的逻辑
func writeDtmResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// 如果是业务定义的失败 (如库存不足), 返回 FAILURE
		if err == dtmcli.ErrFailure {
			w.WriteHeader(http.StatusConflict) // 409 Conflict 更符合语义
			json.NewEncoder(w).Encode(dtmcli.ResultFailure)
		} else {
			// 如果是系统异常, 返回 Ongoing (HTTP 500), DTM 会重试
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		}
	} else {
		// 成功，返回 SUCCESS
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dtmcli.ResultSuccess)
	}
}

func (h *HttpHandler) StockDeductTry(w http.ResponseWriter, r *http.Request) {
	// Try 阶段通常只做资源检查和预留，这里我们简化处理，直接返回成功
	// 在实际项目中，这里可以检查库存是否可能足够
	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.stockSrv.TryDeduct(r.Context(), req)

	writeDtmResponse(w, nil)
}

func (h *HttpHandler) StockDeductConfirm(w http.ResponseWriter, r *http.Request) {
	//logger.Ctx(r.Context()).Println("✅ StockDeductConfirm")
	//
	//var req application.StockActionRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	logger.Ctx(r.Context()).Err(err).Send()
	//	// 请求体解析失败是系统问题，返回错误让 DTM 重试
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//// 调用业务逻辑
	//err := h.lotteryService.DeductStock(r.Context(), &req)
	//logger.Ctx(r.Context()).Err(err).Send()

	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.stockSrv.ConfirmDeduct(r.Context(), req)

	// 使用统一的响应函数
	writeDtmResponse(w, nil)
}

func (h *HttpHandler) StockDeductCancel(w http.ResponseWriter, r *http.Request) {
	//logger.Ctx(r.Context()).Println("✅ StockDeductCancel")
	//
	//var req application.StockActionRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//// 调用业务逻辑
	//err := h.lotteryService.IncreaseStock(r.Context(), &req)
	//
	//// Cancel 阶段必须成功，即使业务逻辑出错，也要返回成功给 DTM
	//// 但需要记录严重错误日志以进行人工干预
	//if err != nil {
	//	logger.Ctx(r.Context()).Err(err).Msg("CRITICAL: StockDeductCancel failed, manual intervention required.")
	//}

	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.stockSrv.CancelDeduct(r.Context(), req)

	writeDtmResponse(w, nil) // 始终向 DTM 返回成功
}

func (h *HttpHandler) AssetTry(w http.ResponseWriter, r *http.Request) {
	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.assetSrv.TryDeduct(r.Context(), req)

	// 假设资产检查成功
	writeDtmResponse(w, nil)
}

func (h *HttpHandler) AssetConfirm(w http.ResponseWriter, r *http.Request) {
	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.assetSrv.ConfirmDeduct(r.Context(), req)

	// 假设资产扣减成功
	writeDtmResponse(w, nil)
}

func (h *HttpHandler) AssetCancel(w http.ResponseWriter, r *http.Request) {
	var req port2.StockActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	h.assetSrv.CancelDeduct(r.Context(), req)

	// 假设资产返还成功
	writeDtmResponse(w, nil)
}
