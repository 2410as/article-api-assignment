// DBに関わる処理
package repository

import (
	"intern-article-api/internal/model"

	"golang.org/x/tools/go/analysis/passes/errorsas"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func (r *ArticleRepository) Save(article *model.Article) error {
	err := r.db.Save(article).Error
	return err
}

func (r *ArticleRepository) FindAll() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Find(&articles).Error
	return articles, err
}