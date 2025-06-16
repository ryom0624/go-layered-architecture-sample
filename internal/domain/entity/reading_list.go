package entity

import (
	"time"
)

type ReadingList struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User             User                `json:"user" gorm:"foreignKey:UserID"`
	ReadingListItems []ReadingListItem   `json:"items,omitempty" gorm:"foreignKey:ReadingListID"`
}

type ReadingListItem struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ReadingListID uint      `json:"reading_list_id" gorm:"not null;index;uniqueIndex:idx_reading_list_article"`
	ArticleID     uint      `json:"article_id" gorm:"not null;index;uniqueIndex:idx_reading_list_article"`
	AddedAt       time.Time `json:"added_at"`
	Notes         string    `json:"notes"`

	// Relations
	ReadingList ReadingList `json:"reading_list" gorm:"foreignKey:ReadingListID"`
	Article     Article     `json:"article" gorm:"foreignKey:ArticleID"`
}