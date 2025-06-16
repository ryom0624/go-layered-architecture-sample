package entity

import "time"

type ArticleView struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ArticleID    uint      `json:"article_id" gorm:"not null;index"`
	UserID       *uint     `json:"user_id" gorm:"index"`
	IPAddress    string    `json:"ip_address" gorm:"not null;index"`
	UserAgent    string    `json:"user_agent"`
	ReadingTime  int       `json:"reading_time" gorm:"default:0"`
	ViewedAt     time.Time `json:"viewed_at" gorm:"not null;index"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Article *Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}