package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"intern-article-api/internal/model"
	"intern-article-api/internal/repository"
)

// FB: externalArticle
type externalArticle struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Medias []struct {
		ID          int    `json:"id"`
		ContentUrl  string `json:"contentUrl"`
		ContentType string `json:"contentType"`
	}
	PublishedAt time.Time `json:"publishedAt"`
}

type ArticleService struct {
	articleRepo      *repository.ArticleRepository
	articleMediaRepo *repository.ArticleMediaRepository
	apiURL           string
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	articleMediaRepo *repository.ArticleMediaRepository,
	apiURL string,
) *ArticleService {
	return &ArticleService{
		articleRepo:      articleRepo,
		articleMediaRepo: articleMediaRepo,
		apiURL:           apiURL,
	}
}

func (s *ArticleService) ImportArticles() error {
	resp, err := http.Get(s.apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("IMPORT status:", resp.Status)
	fmt.Println("IMPORT body head:", string(b[:min(300, len(b))]))

	var externals []externalArticle
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&externals); err != nil {
		return err
	}

	for _, e := range externals {
		article := model.Article{
			ID:          e.ID,
			Title:       e.Title,
			Body:        e.Body,
			PublishedAt: e.PublishedAt,
		}
		if err := s.articleRepo.Save(&article); err != nil {
			return err
		}
		// FB: Media
		for _, m := range e.Medias {
			media := model.ArticleMedia{
				ID:          m.ID,
				ArticleID:   article.ID,
				ContentUrl:  m.ContentUrl,
				ContentType: m.ContentType,
			}
			if err := s.articleMediaRepo.Save(&media); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ArticleService) GetArticles() ([]model.Article, error) {
	articles, err := s.articleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	for i := range articles {
		articleMedias, err := s.articleMediaRepo.FindByArticleID(articles[i].ID)
		if err != nil {
			return nil, err
		}
		articles[i].Medias = articleMedias
	}
	return articles, nil
}

// 入力用構造体（フロント用）
type SaveArticleInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *ArticleService) SaveArticle(input *SaveArticleInput) error {
	article := model.Article{
		Title: input.Title,
		Body:  input.Body,
	}
	return s.articleRepo.Save(&article)
}

// 削除
func (s *ArticleService) DeleteArticle(id int) error {
	if err := s.articleMediaRepo.DeleteByArticleID(id); err != nil {
		return err
	}
	return s.articleRepo.Delete(id)
}

// ピン留め切り替え
func (s *ArticleService) TogglePin(id int) error {
	// FB: Get
	article, err := s.articleRepo.Get(id)
	if err != nil {
		return err
	}
	article.IsPinned = !article.IsPinned
	return s.articleRepo.Save(article)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
