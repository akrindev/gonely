package auth

import "context"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByHandle(ctx context.Context, handle string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id string) (*Account, error)
	FindByUserID(ctx context.Context, userID string) ([]*Account, error)
	FindByProviderAndAccountID(ctx context.Context, providerID, accountID string) (*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
}

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	FindByToken(ctx context.Context, token string) (*Session, error)
	FindByUserID(ctx context.Context, userID string) ([]*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	RevokeAllByUserID(ctx context.Context, userID string) error
}

// VerificationRepository defines the interface for verification data operations
type VerificationRepository interface {
	Create(ctx context.Context, verification *Verification) error
	FindByIdentifier(ctx context.Context, identifier string) (*Verification, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}
