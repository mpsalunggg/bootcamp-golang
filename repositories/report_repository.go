package repositories

import (
	"bootcamp-golang/models"
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReportByDateRange(startDate, endDate time.Time) (*models.ReportResponse, error) {
	start := startDate.Format("2006-01-02")
	end := endDate.Format("2006-01-02")

	resp := &models.ReportResponse{}

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM "transaction"
		WHERE created_at::date >= $1 AND created_at::date <= $2
	`, start, end).Scan(&resp.TotalRevenue, &resp.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	var nama string
	var qty int
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_detail td
		JOIN "transaction" t ON t.id = td.transaction_id
		JOIN product p ON p.id = td.product_id
		WHERE t.created_at::date >= $1 AND t.created_at::date <= $2
		GROUP BY td.product_id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, start, end).Scan(&nama, &qty)
	if err == nil {
		resp.ProdukTerlaris = &models.ProdukTerlaris{Nama: nama, QuantityTerjual: qty}
	}

	return resp, nil
}
