package auth

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session entity
type Session struct {
	ID            string     `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID        string     `json:"userId" gorm:"type:varchar(255);not null;index:session_user_id_idx"`
	Token         string     `json:"token" gorm:"type:varchar(500);not null;uniqueIndex:session_token_idx"`
	UserAgent     *string    `json:"userAgent,omitempty" gorm:"type:varchar(500)"`
	IPAddress     *string    `json:"ipAddress,omitempty" gorm:"type:varchar(45)"`
	Revoked       bool       `json:"revoked" gorm:"not null;default:false"`
	RevokedAt     *time.Time `json:"revokedAt,omitempty"`
	RevokedReason *string    `json:"revokedReason,omitempty" gorm:"type:varchar(255)"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt    *time.Time `json:"lastUsedAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"not null;index:session_created_at_idx"`
	UpdatedAt     time.Time  `json:"updatedAt" gorm:"not null"`

	// Relation
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// NewSession creates a new session
func NewSession(userID, token string, expiresAt time.Time, userAgent, ipAddress *string) *Session {
	now := time.Now()
	return &Session{
		ID:         uuid.New().String(),
		UserID:     userID,
		Token:      token,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
		Revoked:    false,
		ExpiresAt:  &expiresAt,
		LastUsedAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsValid checks if the session is valid (not revoked and not expired)
func (s *Session) IsValid() bool {
	if s.Revoked {
		return false
	}
	if s.ExpiresAt != nil && time.Now().After(*s.ExpiresAt) {
		return false
	}
	return true
}

// Revoke revokes the session with a reason
func (s *Session) Revoke(reason string) {
	now := time.Now()
	s.Revoked = true
	s.RevokedAt = &now
	s.RevokedReason = &reason
	s.UpdatedAt = now
}

// UpdateLastUsed updates the last used timestamp
func (s *Session) UpdateLastUsed() {
	now := time.Now()
	s.LastUsedAt = &now
	s.UpdatedAt = now
}

// Rotate generates a new token for the session
func (s *Session) Rotate(newToken string, newExpiresAt time.Time) {
	s.Token = newToken
	s.ExpiresAt = &newExpiresAt
	s.UpdatedAt = time.Now()
}
