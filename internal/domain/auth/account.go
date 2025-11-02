package auth

import (
	"time"

	"github.com/google/uuid"
)

// Account represents an authentication account entity
type Account struct {
	ID                     string     `json:"id" gorm:"primaryKey;type:varchar(255)"`
	AccountID              string     `json:"accountId" gorm:"type:varchar(255);not null"`
	ProviderID             string     `json:"providerId" gorm:"type:varchar(255);not null"`
	UserID                 string     `json:"userId" gorm:"type:varchar(255);not null;index:account_user_id_idx"`
	AccessToken            *string    `json:"accessToken,omitempty" gorm:"type:varchar(500);index:account_access_token_idx"`
	AccessTokenExpiresAt   *time.Time `json:"accessTokenExpiresAt,omitempty"`
	RefreshToken           *string    `json:"refreshToken,omitempty" gorm:"type:varchar(500)"`
	RefreshTokenExpiresAt  *time.Time `json:"refreshTokenExpiresAt,omitempty"`
	IDToken                *string    `json:"idToken,omitempty" gorm:"type:text"`
	Scope                  *string    `json:"scope,omitempty" gorm:"type:varchar(255)"`
	Password               *string    `json:"-" gorm:"type:varchar(255)"` // Never expose in JSON
	CreatedAt              time.Time  `json:"createdAt" gorm:"not null;index:account_created_at_idx"`
	UpdatedAt              time.Time  `json:"updatedAt" gorm:"not null"`

	// Relation
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for Account
func (Account) TableName() string {
	return "accounts"
}

// NewAccount creates a new account
func NewAccount(accountID, providerID, userID string) *Account {
	now := time.Now()
	return &Account{
		ID:         uuid.New().String(),
		AccountID:  accountID,
		ProviderID: providerID,
		UserID:     userID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// SetTokens sets OAuth tokens for the account
func (a *Account) SetTokens(accessToken, refreshToken string, accessExpiry, refreshExpiry *time.Time) {
	a.AccessToken = &accessToken
	a.RefreshToken = &refreshToken
	a.AccessTokenExpiresAt = accessExpiry
	a.RefreshTokenExpiresAt = refreshExpiry
	a.UpdatedAt = time.Now()
}

// SetPassword sets the hashed password for the account
func (a *Account) SetPassword(hashedPassword string) {
	a.Password = &hashedPassword
	a.UpdatedAt = time.Now()
}

// IsTokenValid checks if the access token is still valid
func (a *Account) IsTokenValid() bool {
	if a.AccessToken == nil || a.AccessTokenExpiresAt == nil {
		return false
	}
	return time.Now().Before(*a.AccessTokenExpiresAt)
}
