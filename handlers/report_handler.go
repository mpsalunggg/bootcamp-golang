package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bootcamp-golang/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetReportHariIni(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	report, err := h.service.GetReportHariIni()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetReportByDateRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	startStr := q.Get("start_date")
	endStr := q.Get("end_date")
	if startStr == "" || endStr == "" {
		http.Error(w, "start_date dan end_date harus diisi", http.StatusBadRequest)
		return
	}
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, "format start_date tidak valid (gunakan YYYY-MM-DD)", http.StatusBadRequest)
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		http.Error(w, "format end_date tidak valid (gunakan YYYY-MM-DD)", http.StatusBadRequest)
		return
	}
	if end.Before(start) {
		http.Error(w, "end_date harus >= start_date", http.StatusBadRequest)
		return
	}
	report, err := h.service.GetReportByDateRange(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
