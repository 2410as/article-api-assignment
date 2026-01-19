package repository

import (
	"intern-article-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type ArticleMediaRepository struct {
	db *gorm.DB
}

func NewArticleMediaRepository(db *gorm.DB) *ArticleMediaRepository {
	return &ArticleMediaRepository{db: db}
}

func (r *ArticleMediaRepository) Save(articleMedia *model.ArticleMedia) error {
	articleMedia.UpdatedAt = time.Now()
	if articleMedia.CreatedAt.IsZero() {
		articleMedia.CreatedAt = articleMedia.UpdatedAt
	}
	return r.db.Save(articleMedia).Error
}

func (r *ArticleMediaRepository) FindByArticleID(articleID int) ([]model.ArticleMedia, error) {
	var articleMedias []model.ArticleMedia
	err := r.db.Find(&articleMedias, "article_id = ?", articleID).Error
	return articleMedias, err
}

func (r *ArticleMediaRepository) DeleteByArticleID(articleID int) error {
	return r.db.Delete(&model.ArticleMedia{}, "article_id = ?", articleID).Error
}
