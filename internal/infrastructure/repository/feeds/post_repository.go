package feeds

import (
	"context"

	"github.com/akrindev/gonely/internal/domain/feeds"
	"gorm.io/gorm"
)

// postRepository implements the feeds.PostRepository interface
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) feeds.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *feeds.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) FindByID(ctx context.Context, id string) (*feeds.Post, error) {
	var post feeds.Post
	err := r.db.WithContext(ctx).
		Preload("Comments").
		Where("id = ?", id).
		First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll(ctx context.Context, limit, offset int) ([]*feeds.Post, error) {
	var posts []*feeds.Post
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(ctx context.Context, post *feeds.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&feeds.Post{}, "id = ?", id).Error
}
