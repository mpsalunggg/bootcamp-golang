package main

import (
	"bootcamp-golang/database"
	"bootcamp-golang/handlers"
	"bootcamp-golang/repositories"
	"bootcamp-golang/services"
	"bootcamp-golang/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

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
	categoryRepository := repositories.NewCategoryRepository(db)
	productService := services.NewProductService(productRepository, categoryRepository)
	categoryService := services.NewCategoryService(categoryRepository)
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	_ = repositories.NewCategoryRepository(db)

	http.HandleFunc("/api/produk", productHandler.HandleProduct)
	http.HandleFunc("/api/produk/", productHandler.HandleProductById)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryById)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.NewResponse("success", "OK"))
	})

	log.Println("Server running di localhost:" + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("gagal running server:", err)
	}
}
