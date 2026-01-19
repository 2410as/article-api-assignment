package model

import "time"

// FB: Model
type ArticleMedia struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	ArticleID   int       `gorm:"index" json:"article_id"`
	ContentUrl  string    `json:"content_url"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
