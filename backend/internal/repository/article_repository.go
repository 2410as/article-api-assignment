package repository

import (
	"intern-article-api/internal/model"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Save(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) FindAll() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) Delete(id int) error {
	return r.db.Delete(&model.Article{}, id).Error
}
