package entity

import "time"

type ArticleStatistics struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	ArticleID            uint      `json:"article_id" gorm:"not null;uniqueIndex"`
	TotalViews           int       `json:"total_views" gorm:"default:0"`
	UniqueViews          int       `json:"unique_views" gorm:"default:0"`
	AuthenticatedViews   int       `json:"authenticated_views" gorm:"default:0"`
	AnonymousViews       int       `json:"anonymous_views" gorm:"default:0"`
	AverageReadingTime   float32   `json:"average_reading_time" gorm:"default:0"`
	CompletionRate       float32   `json:"completion_rate" gorm:"default:0"`
	TotalReadingTime     int       `json:"total_reading_time" gorm:"default:0"`
	LastCalculatedAt     time.Time `json:"last_calculated_at" gorm:"not null"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relations
	Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}