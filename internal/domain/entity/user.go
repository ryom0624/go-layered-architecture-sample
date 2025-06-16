package entity

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Favorites         []Favorite            `json:"favorites,omitempty" gorm:"foreignKey:UserID"`
	ReadingLists      []ReadingList         `json:"reading_lists,omitempty" gorm:"foreignKey:UserID"`
	Views             []ArticleView         `json:"views,omitempty" gorm:"foreignKey:UserID"`
	ReadingHistories  []UserReadingHistory  `json:"reading_histories,omitempty" gorm:"foreignKey:UserID"`
}