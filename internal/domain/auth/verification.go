package auth

import (
	"time"

	"github.com/google/uuid"
)

// Verification represents an email/identity verification entity
type Verification struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Identifier string    `json:"identifier" gorm:"type:varchar(255);not null;index:verification_identifier_idx"`
	Value      string    `json:"value" gorm:"type:varchar(500);not null"`
	ExpiresAt  time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"not null"`
}

// TableName specifies the table name for Verification
func (Verification) TableName() string {
	return "verifications"
}

// NewVerification creates a new verification
func NewVerification(identifier, value string, expiresAt time.Time) *Verification {
	now := time.Now()
	return &Verification{
		ID:         uuid.New().String(),
		Identifier: identifier,
		Value:      value,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsExpired checks if the verification has expired
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsValid checks if the verification is valid (not expired)
func (v *Verification) IsValid() bool {
	return !v.IsExpired()
}
