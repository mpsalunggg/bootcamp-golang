package services

import (
	"bootcamp-golang/models"
	"bootcamp-golang/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReportHariIni() (*models.ReportResponse, error) {
	now := time.Now()
	return s.repo.GetReportByDateRange(now, now)
}

func (s *ReportService) GetReportByDateRange(startDate, endDate time.Time) (*models.ReportResponse, error) {
	return s.repo.GetReportByDateRange(startDate, endDate)
}
