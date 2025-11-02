package feeds

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment entity
type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	PostID    string    `json:"postId" gorm:"type:varchar(255);not null;index:comment_post_id_idx"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;index:comment_created_at_idx"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`

	// Relation
	Post *Post `json:"-" gorm:"foreignKey:PostID"`
}

// TableName specifies the table name for Comment
func (Comment) TableName() string {
	return "comments"
}

// NewComment creates a new comment
func NewComment(postID, content string) *Comment {
	now := time.Now()
	return &Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateContent updates the comment content
func (c *Comment) UpdateContent(content string) {
	c.Content = content
	c.UpdatedAt = time.Now()
}
