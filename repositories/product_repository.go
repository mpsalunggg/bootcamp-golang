package repositories

import (
	"bootcamp-golang/models"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name as category_name FROM product p LEFT JOIN category c ON p.category_id = c.id"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var p models.Product
		var categoryName sql.NullString
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryName)
		if err != nil {
			return nil, err
		}
		if categoryName.Valid {
			p.Category = &models.Category{ID: p.CategoryID, Name: categoryName.String}
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	return repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}
