package feeds

import (
	"time"

	"github.com/google/uuid"
)

// Post represents the post aggregate root in the feeds domain
type Post struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null;index:post_created_at_idx"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`

	// Relations
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Post
func (Post) TableName() string {
	return "posts"
}

// NewPost creates a new post with the given content
func NewPost(content string) *Post {
	now := time.Now()
	return &Post{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateContent updates the post content
func (p *Post) UpdateContent(content string) {
	p.Content = content
	p.UpdatedAt = time.Now()
}

// AddComment adds a comment to the post
func (p *Post) AddComment(comment *Comment) {
	p.Comments = append(p.Comments, *comment)
}
