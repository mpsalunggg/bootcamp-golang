package services

import (
	"bootcamp-golang/models"
	"bootcamp-golang/repositories"
	"errors"
)

type ProductService struct {
	repo         *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(product *models.Product) error {
	category, err := s.categoryRepo.GetById(*product.CategoryID)
	if err != nil {
		return err
	}

	if category == nil {
		return errors.New("Kategori tidak ditemukan")
	}

	return s.repo.Create(product)
}

func (s *ProductService) GetById(id int) (*models.Product, error) {
	return s.repo.GetById(id)
}

func (s *ProductService) Update(product *models.Product) error {
	category, err := s.categoryRepo.GetById(*product.CategoryID)
	if err != nil {
		return err
	}

	if category == nil {
		return errors.New("Kategori tidak ditemukan")
	}
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
