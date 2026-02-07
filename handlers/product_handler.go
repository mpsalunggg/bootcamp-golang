package handlers

import (
	"bootcamp-golang/models"
	"bootcamp-golang/services"
	"bootcamp-golang/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func parseProductID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/api/produk/")
	idStr = strings.Trim(idStr, "/")
	if idx := strings.Index(idStr, "/"); idx > 0 {
		idStr = idStr[:idx]
	}
	return strconv.Atoi(idStr)
}

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetById(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	products, err := h.service.GetAll(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.NewResponse("Berhasil mendapatkan semua produk", products))
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
		return
	}
	err = h.service.Create(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.NewResponse("success create product", product))
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}
	product, err := h.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.NewResponse("Berhasil mendapatkan produk berdasarkan id", product))
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}
	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.Update(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewResponse("Berhasil memperbarui produk", product))
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewResponse("Berhasil menghapus produk", nil))
}
