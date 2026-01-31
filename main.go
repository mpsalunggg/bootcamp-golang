package main

import (
	"bootcamp-golang/database"
	"bootcamp-golang/handlers"
	"bootcamp-golang/repositories"
	"bootcamp-golang/services"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

// const port = ":8080"

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var categories = []Category{}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	db, err := database.InitDB(config.DatabaseURL)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/produk", productHandler.HandleProduct)

	// http.HandleFunc("/categories", handleCategories)
	// http.HandleFunc("/categories/", handleCategoryById)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Message: "success",
			Data:    "OK",
		})
	})

	log.Println("Server running di localhost:" + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("gagal running server:", err)
	}

	// fmt.Println("server is running on port " + port)

	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	fmt.Println("running server error: ", err)
	// }
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getAllCategories(w)
		return
	}
	if r.Method == "POST" {
		createCategory(w, r)
		return
	}
}

func handleCategoryById(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getCategoryById(w, r)
		return
	}
	if r.Method == "PUT" {
		updateCategory(w, r)
		return
	}
	if r.Method == "DELETE" {
		deleteCategory(w, r)
		return
	}
}

func getAllCategories(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "success get all data",
		Data:    categories,
	})
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var category Category
	json.NewDecoder(r.Body).Decode(&category)
	category.ID = len(categories) + 1
	categories = append(categories, category)
	json.NewEncoder(w).Encode(Response{
		Message: "success create data",
		Data:    category,
	})
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid ID"})
		return
	}

	for _, cat := range categories {
		if cat.ID == id {
			json.NewEncoder(w).Encode(Response{
				Message: "success get data",
				Data:    cat,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Category not found"})
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid ID"})
		return
	}

	var updatedCat Category
	json.NewDecoder(r.Body).Decode(&updatedCat)

	for i, cat := range categories {
		if cat.ID == id {
			updatedCat.ID = id
			categories[i] = updatedCat
			json.NewEncoder(w).Encode(Response{
				Message: "success update data",
				Data:    updatedCat,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Category not found"})
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "invalid ID"})
		return
	}

	for i, cat := range categories {
		if cat.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			json.NewEncoder(w).Encode(Response{
				Message: "success delete data",
				Data:    cat,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Category not found"})
}
