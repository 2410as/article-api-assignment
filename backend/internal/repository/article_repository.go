package repository

import (
	"intern-article-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Save(article *model.Article) error {
	article.UpdatedAt = time.Now()
	if article.CreatedAt.IsZero() {
		article.CreatedAt = article.UpdatedAt
	}
	return r.db.Save(article).Error
}

func (r *ArticleRepository) FindAll() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Find(&articles).Error
	return articles, err
}

// FB: Get
func (r *ArticleRepository) Get(id int) (*model.Article, error) {
	var article model.Article
	err := r.db.First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) Delete(id int) error {
	return r.db.Delete(&model.Article{}, id).Error
}
