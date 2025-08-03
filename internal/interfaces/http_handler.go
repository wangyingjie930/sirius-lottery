package interfaces

import (
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
	mux.HandleFunc("POST /api/v2/lottery/draw", h.Draw)                               // [cite: 154]
	mux.HandleFunc("GET /api/v2/lottery/instance/{instanceId}", h.GetLotteryInstance) // [cite: 167]
}

func (h *HttpHandler) Draw(writer http.ResponseWriter, request *http.Request) {
	var req application.DrawRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.lotteryService.Draw(request.Context(), &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

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

	json.NewEncoder(writer).Encode(resp)
}
