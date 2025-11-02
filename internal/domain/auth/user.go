package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user aggregate root in the authentication domain
type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	Handle        *string   `json:"handle,omitempty" gorm:"type:varchar(255);index:user_handle_idx"`
	Email         string    `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	EmailVerified bool      `json:"emailVerified" gorm:"not null;default:false"`
	IsAnonymous   bool      `json:"isAnonymous" gorm:"default:false"`
	Image         *string   `json:"image,omitempty" gorm:"type:varchar(500)"`
	CreatedAt     time.Time `json:"createdAt" gorm:"not null;index:user_created_at_idx"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"not null"`

	// Relations
	Accounts []Account `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Sessions []Session `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with default values
func NewUser(name, email string, isAnonymous bool) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New().String(),
		Name:          name,
		Email:         email,
		EmailVerified: false,
		IsAnonymous:   isAnonymous,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// SetHandle sets the user's handle
func (u *User) SetHandle(handle string) {
	u.Handle = &handle
	u.UpdatedAt = time.Now()
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}

// SetImage sets the user's profile image
func (u *User) SetImage(imageURL string) {
	u.Image = &imageURL
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates user's profile information
func (u *User) UpdateProfile(name string, handle *string) {
	u.Name = name
	u.Handle = handle
	u.UpdatedAt = time.Now()
}
