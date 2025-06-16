package entity

import "time"

type DailyStatistics struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	Date                  time.Time `json:"date" gorm:"not null;uniqueIndex;type:date"`
	TotalViews            int       `json:"total_views" gorm:"default:0"`
	UniqueUsers           int       `json:"unique_users" gorm:"default:0"`
	NewUsers              int       `json:"new_users" gorm:"default:0"`
	TotalReadingTime      int       `json:"total_reading_time" gorm:"default:0"`
	AverageReadingTime    float32   `json:"average_reading_time" gorm:"default:0"`
	PopularArticleID      *uint     `json:"popular_article_id"`
	TopReaderUserID       *uint     `json:"top_reader_user_id"`
	TotalArticles         int       `json:"total_articles" gorm:"default:0"`
	NewArticles           int       `json:"new_articles" gorm:"default:0"`
	PublishedArticles     int       `json:"published_articles" gorm:"default:0"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Relations
	PopularArticle *Article `json:"popular_article,omitempty" gorm:"foreignKey:PopularArticleID"`
	TopReaderUser  *User    `json:"top_reader_user,omitempty" gorm:"foreignKey:TopReaderUserID"`
}