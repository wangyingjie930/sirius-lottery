package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"sirius-lottery/internal/application"
)

type HttpHandler struct {
	lotteryService application.LotteryService
}

func NewHttpHandler(lotteryService application.LotteryService) *HttpHandler {
	return &HttpHandler{lotteryService: lotteryService}
}

func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v2/lottery/draw", h.Draw)                               //
	mux.HandleFunc("GET /api/v2/lottery/instance/{instanceId}", h.GetLotteryInstance) //
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

