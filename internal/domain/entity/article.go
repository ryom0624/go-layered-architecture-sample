package entity

import "time"

type Article struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"not null"`
	Content       string    `json:"content" gorm:"type:text"`
	AuthorID      uint      `json:"author_id" gorm:"not null"`
	Author        User      `json:"author" gorm:"foreignKey:AuthorID"`
	Status        string    `json:"status" gorm:"default:'draft'"`
	FavoriteCount int       `json:"favorite_count" gorm:"default:0"`
	ViewCount     int       `json:"view_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Favorites             []Favorite             `json:"favorites,omitempty" gorm:"foreignKey:ArticleID"`
	ReadingListItems      []ReadingListItem      `json:"reading_list_items,omitempty" gorm:"foreignKey:ArticleID"`
	Views                 []ArticleView          `json:"views,omitempty" gorm:"foreignKey:ArticleID"`
	ReadingHistories      []UserReadingHistory   `json:"reading_histories,omitempty" gorm:"foreignKey:ArticleID"`
	Statistics            *ArticleStatistics     `json:"statistics,omitempty" gorm:"foreignKey:ArticleID"`
}