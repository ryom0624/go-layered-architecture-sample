package entity

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	Slug      string    `json:"slug" gorm:"unique;not null"`
	Color     string    `json:"color" gorm:"default:'#3B82F6'"` // Default blue color
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`
}

// ValidateName validates tag name
func (t *Tag) ValidateName() error {
	if t.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if utf8.RuneCountInString(t.Name) < 2 {
		return &ValidationError{Field: "name", Message: "name must be at least 2 characters"}
	}
	if utf8.RuneCountInString(t.Name) > 50 {
		return &ValidationError{Field: "name", Message: "name must be 50 characters or less"}
	}
	return nil
}

// ValidateColor validates tag color (hex color code)
func (t *Tag) ValidateColor() error {
	if t.Color == "" {
		t.Color = "#3B82F6" // Default blue color
		return nil
	}
	
	// Check if it's a valid hex color code
	matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, t.Color)
	if !matched {
		return &ValidationError{Field: "color", Message: "color must be a valid hex color code (e.g., #FF0000)"}
	}
	return nil
}

// GenerateSlug generates URL slug from tag name
func (t *Tag) GenerateSlug() {
	if t.Name == "" {
		return
	}
	t.Slug = generateTagSlug(t.Name)
}

// Validate validates all tag fields
func (t *Tag) Validate() error {
	if err := t.ValidateName(); err != nil {
		return err
	}
	if err := t.ValidateColor(); err != nil {
		return err
	}
	return nil
}

// generateTagSlug generates slug from tag name
func generateTagSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	
	// Japanese and common tag mappings
	jpMappings := map[string]string{
		"Go言語":         "golang",
		"Python":       "python",
		"JavaScript":   "javascript",
		"React":        "react",
		"Vue.js":       "vuejs",
		"Node.js":      "nodejs",
		"Docker":       "docker",
		"Kubernetes":   "kubernetes",
		"AWS":          "aws",
		"Azure":        "azure",
		"GCP":          "gcp",
		"PostgreSQL":   "postgresql",
		"MySQL":        "mysql",
		"MongoDB":      "mongodb",
		"Redis":        "redis",
		"Git":          "git",
		"GitHub":       "github",
		"GitLab":       "gitlab",
		"CI/CD":        "cicd",
		"API":          "api",
		"REST":         "rest",
		"GraphQL":      "graphql",
		"JSON":         "json",
		"XML":          "xml",
		"HTML":         "html",
		"CSS":          "css",
		"TypeScript":   "typescript",
		"Java":         "java",
		"C++":          "cpp",
		"C#":           "csharp",
		"Ruby":         "ruby",
		"PHP":          "php",
		"Swift":        "swift",
		"Kotlin":       "kotlin",
		"Rust":         "rust",
		"初心者":        "beginner",
		"入門":          "introduction",
		"基礎":          "basics",
		"応用":          "advanced",
		"チュートリアル": "tutorial",
		"ガイド":        "guide",
		"実践":          "practice",
		"サンプル":      "sample",
		"例":           "example",
	}
	
	// Check for mappings
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
		slug = "tag"
	}
	
	return slug
}

// TagWithCount represents tag with usage count
type TagWithCount struct {
	*Tag
	UsageCount int `json:"usage_count"`
}