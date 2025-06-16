package entity

import (
	"time"
	"layered-architecture-template/pkg/constants"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	ParentID  *uint     `json:"parent_id,omitempty" gorm:"index"` // 返信の場合の親コメントID
	Status    string    `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Author   User      `json:"author" gorm:"foreignKey:AuthorID"`
	Article  Article   `json:"article" gorm:"foreignKey:ArticleID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// ValidateContent validates comment content
func (c *Comment) ValidateContent() error {
	if c.Content == "" {
		return &ValidationError{Field: "content", Message: "content is required"}
	}
	if len(c.Content) > constants.MaxCommentLength {
		return &ValidationError{Field: "content", Message: "content must be 1000 characters or less"}
	}
	return nil
}

// ValidateStatus validates comment status
func (c *Comment) ValidateStatus() error {
	validStatuses := []string{"pending", "approved", "rejected"}
	for _, status := range validStatuses {
		if c.Status == status {
			return nil
		}
	}
	return &ValidationError{Field: "status", Message: "status must be one of: pending, approved, rejected"}
}

// IsApproved returns true if comment is approved
func (c *Comment) IsApproved() bool {
	return c.Status == "approved"
}

// IsPending returns true if comment is pending approval
func (c *Comment) IsPending() bool {
	return c.Status == "pending"
}

// IsReply returns true if this is a reply to another comment
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}