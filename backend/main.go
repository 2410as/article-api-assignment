package main

import (
	"log"
	"net/http"
	"os"

	"intern-article-api/internal/handler"
	"intern-article-api/internal/model"
	"intern-article-api/internal/repository"
	"intern-article-api/internal/service"

	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	apiURL := os.Getenv("EXTERNAL_API_URL")
	if apiURL == "" {
		log.Fatal("EXTERNAL_API_URL is not set")
	}

	wd, _ := os.Getwd()
	dbPath := wd + "/article.db"
	log.Println("DB PATH:", dbPath)
	

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.Article{}); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewArticleRepository(db)
	svc := service.NewArticleService(repo, apiURL)
	h := handler.NewArticleHandler(svc)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	
	r.Post("/articles/import", h.ImportArticles)
	r.Get("/articles", h.GetArticles)
	r.Post("/articles", h.CreateArticle)
	r.Patch("/articles/{id}/pin", h.TogglePinArticle)
	r.Delete("/articles/{id}", h.DeleteArticle) // ← 追加

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}
