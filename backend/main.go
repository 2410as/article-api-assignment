package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"intern-article-api/internal/handler"
	"intern-article-api/internal/model"
	"intern-article-api/internal/repository"
	"intern-article-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found; using environment variables")
	}

	apiURL := os.Getenv("EXTERNAL_API_URL")
	if apiURL == "" {
		log.Fatal("EXTERNAL_API_URL is not set")
	}

	dbPath := "article.db"
	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB PATH:", absDBPath)

	db, err := gorm.Open(sqlite.Open(absDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.Article{}, &model.ArticleMedia{}); err != nil {
		log.Fatal(err)
	}

	articleRepo := repository.NewArticleRepository(db)
	articleMediaRepo := repository.NewArticleMediaRepository(db)
	svc := service.NewArticleService(articleRepo, articleMediaRepo, apiURL)
	h := handler.NewArticleHandler(svc)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// FB: Route
	r.Route("/articles", func(r chi.Router) {
		r.Post("/import", h.ImportArticles)
		r.Get("/", h.GetArticles)
		r.Post("/", h.CreateArticle)
		r.Patch("/{id}/pin", h.TogglePinArticle)
		r.Delete("/{id}", h.DeleteArticle)
	})

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
