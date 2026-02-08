package repositories

import (
	"bootcamp-golang/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name FROM product p LEFT JOIN category c ON p.category_id = c.id"
	args := []interface{}{}

	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var p models.Product
		var categoryName sql.NullString
		var categoryID sql.NullInt64
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}
		if categoryID.Valid {
			categoryIDInt := int(categoryID.Int64)
			p.CategoryID = &categoryIDInt
		}
		if categoryName.Valid {
			p.Category = &models.Category{ID: int(categoryID.Int64), Name: categoryName.String}
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	return repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}

func (repo *ProductRepository) GetById(id int) (*models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name FROM product p LEFT JOIN category c ON p.category_id = c.id WHERE p.id = $1"
	var product models.Product
	var categoryName sql.NullString
	var categoryID sql.NullInt64
	err := repo.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &categoryID, &categoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Produk tidak ditemukan")
		}
		return nil, err
	}

	if categoryID.Valid {
		categoryIDInt := int(categoryID.Int64)
		product.CategoryID = &categoryIDInt
	}
	if categoryName.Valid {
		product.Category = &models.Category{ID: int(categoryID.Int64), Name: categoryName.String}
	}

	return &product, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Produk tidak ditemukan")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return nil
}
