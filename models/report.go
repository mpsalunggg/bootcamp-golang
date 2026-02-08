package models

type ProdukTerlaris struct {
	Nama            string `json:"nama"`
	QuantityTerjual int    `json:"jumlah"`
}

type ReportResponse struct {
	TotalRevenue   int             `json:"total_revenue"`
	TotalTransaksi int             `json:"total_transaksi"`
	ProdukTerlaris *ProdukTerlaris `json:"produk_terlaris"`
}
