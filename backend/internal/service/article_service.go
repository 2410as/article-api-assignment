package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"intern-article-api/internal/model"
	"intern-article-api/internal/repository"
)

type externalArticle struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	IsPinned bool   `json:"is_pinned"`
}

type ArticleService struct {
	repo   *repository.ArticleRepository
	apiURL string
}

// 入力用構造体（フロント用）
type ArticleInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewArticleService(repo *repository.ArticleRepository, apiURL string) *ArticleService {
	return &ArticleService{repo: repo, apiURL: apiURL}
}

// 外部記事のインポート
func (s *ArticleService) ImportArticles() error {
	resp, err := http.Get(s.apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var externals []externalArticle
	if err := json.NewDecoder(resp.Body).Decode(&externals); err != nil {
		return err
	}

	for _, e := range externals {
		article := model.Article{
			ID:       e.ID,
			Title:    e.Title,
			Body:     e.Body,
			IsPinned: e.IsPinned,
		}
		if err := s.repo.Save(&article); err != nil {
			return err
		}
	}
	return nil
}

func (s *ArticleService) GetArticles() ([]model.Article, error) {
	return s.repo.FindAll()
}

func (s *ArticleService) SaveArticle(input *ArticleInput) error {
	article := model.Article{
		Title: input.Title,
		Body:  input.Body,
	}
	return s.repo.Save(&article)
}

// 削除
func (s *ArticleService) DeleteArticle(id int) error {
	return s.repo.Delete(id)
}

// ピン留め切り替え
func (s *ArticleService) TogglePin(id int) error {
	articles, err := s.repo.FindAll()
	if err != nil {
		return err
	}
	for _, a := range articles {
		if a.ID == id {
			a.IsPinned = !a.IsPinned
			return s.repo.Save(&a)
		}
	}
	return fmt.Errorf("article not found")
}
