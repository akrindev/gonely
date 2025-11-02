package auth

import (
	"context"
	"time"

	"github.com/akrindev/gonely/internal/domain/auth"
	"gorm.io/gorm"
)

// sessionRepository implements the auth.SessionRepository interface
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) auth.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *auth.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) FindByID(ctx context.Context, id string) (*auth.Session, error) {
	var session auth.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*auth.Session, error) {
	var session auth.Session
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ?", token).
		First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID string) ([]*auth.Session, error) {
	var sessions []*auth.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepository) Update(ctx context.Context, session *auth.Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&auth.Session{}, "id = ?", id).Error
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&auth.Session{}).Error
}

func (r *sessionRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	reason := "All sessions revoked"
	return r.db.WithContext(ctx).
		Model(&auth.Session{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":        true,
			"revoked_at":     now,
			"revoked_reason": reason,
			"updated_at":     now,
		}).Error
}
