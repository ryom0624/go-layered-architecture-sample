package entity

import "time"

type UserReadingHistory struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	UserID              uint      `json:"user_id" gorm:"not null;index"`
	ArticleID           uint      `json:"article_id" gorm:"not null;index"`
	TotalReadingTime    int       `json:"total_reading_time" gorm:"default:0"`
	ReadingProgress     float32   `json:"reading_progress" gorm:"default:0"`
	IsCompleted         bool      `json:"is_completed" gorm:"default:false"`
	FirstViewedAt       time.Time `json:"first_viewed_at" gorm:"not null"`
	LastViewedAt        time.Time `json:"last_viewed_at" gorm:"not null"`
	CompletedAt         *time.Time `json:"completed_at"`
	ViewCount           int       `json:"view_count" gorm:"default:1"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relations
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}