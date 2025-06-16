package entity

import (
	"time"
)

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index;uniqueIndex:idx_user_article_favorite"`
	ArticleID uint      `json:"article_id" gorm:"not null;index;uniqueIndex:idx_user_article_favorite"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}