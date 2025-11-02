package feeds

import (
	"context"

	"github.com/akrindev/gonely/internal/domain/feeds"
	"gorm.io/gorm"
)

// commentRepository implements the feeds.CommentRepository interface
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) feeds.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *feeds.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByID(ctx context.Context, id string) (*feeds.Comment, error) {
	var comment feeds.Comment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByPostID(ctx context.Context, postID string, limit, offset int) ([]*feeds.Comment, error) {
	var comments []*feeds.Comment
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Update(ctx context.Context, comment *feeds.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&feeds.Comment{}, "id = ?", id).Error
}
