package auth

import (
	"context"

	"github.com/akrindev/gonely/internal/domain/auth"
	"gorm.io/gorm"
)

// accountRepository implements the auth.AccountRepository interface
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) auth.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *auth.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) FindByID(ctx context.Context, id string) (*auth.Account, error) {
	var account auth.Account
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID string) ([]*auth.Account, error) {
	var accounts []*auth.Account
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindByProviderAndAccountID(ctx context.Context, providerID, accountID string) (*auth.Account, error) {
	var account auth.Account
	err := r.db.WithContext(ctx).
		Where("provider_id = ? AND account_id = ?", providerID, accountID).
		First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Update(ctx context.Context, account *auth.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&auth.Account{}, "id = ?", id).Error
}
