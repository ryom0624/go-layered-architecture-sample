package entity

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}

// ValidateName validates category name
func (c *Category) ValidateName() error {
	if c.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if utf8.RuneCountInString(c.Name) < 2 {
		return &ValidationError{Field: "name", Message: "name must be at least 2 characters"}
	}
	if utf8.RuneCountInString(c.Name) > 100 {
		return &ValidationError{Field: "name", Message: "name must be 100 characters or less"}
	}
	return nil
}

// ValidateDescription validates category description
func (c *Category) ValidateDescription() error {
	if utf8.RuneCountInString(c.Description) > 500 {
		return &ValidationError{Field: "description", Message: "description must be 500 characters or less"}
	}
	return nil
}

// GenerateSlug generates URL slug from category name
func (c *Category) GenerateSlug() {
	if c.Name == "" {
		return
	}
	c.Slug = generateCategorySlug(c.Name)
}

// Validate validates all category fields
func (c *Category) Validate() error {
	if err := c.ValidateName(); err != nil {
		return err
	}
	if err := c.ValidateDescription(); err != nil {
		return err
	}
	return nil
}

// generateCategorySlug generates slug from category name
func generateCategorySlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	
	// Japanese name mapping
	jpMappings := map[string]string{
		"プログラミング":    "programming",
		"ウェブ開発":       "web-development",
		"データベース":     "database",
		"デブオプス":       "devops",
		"機械学習":        "machine-learning",
		"人工知能":        "artificial-intelligence",
		"セキュリティ":     "security",
		"クラウド":        "cloud",
		"モバイル":        "mobile",
		"ゲーム開発":      "game-development",
	}
	
	// Check for Japanese mappings
	if mapped, exists := jpMappings[name]; exists {
		return mapped
	}
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	// If slug is empty after processing, use fallback
	if slug == "" {
		slug = "category"
	}
	
	return slug
}

// CategoryWithCount represents category with article count
type CategoryWithCount struct {
	*Category
	ArticleCount int `json:"article_count"`
}