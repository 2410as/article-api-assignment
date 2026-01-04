package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"intern-article-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(s *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: s,
	}
}


func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")        
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteArticle(id); err != nil { // service層の削除呼び出し
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// GET /articles
func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(articles)
}

// POST /articles
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article service.ArticleInput
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.SaveArticle(&article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// PATCH /articles/{id}/pin
func (h *ArticleHandler) TogglePinArticle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.TogglePin(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// POST /articles/import
func (h *ArticleHandler) ImportArticles(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ImportArticles(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
