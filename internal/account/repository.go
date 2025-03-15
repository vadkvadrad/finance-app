package account

import (
	"finance-app/pkg/db"

	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	Db *db.Db
}

func NewAccountRepository(db *db.Db) *AccountRepository {
	return &AccountRepository{
		Db: db,
	}
}

func (repo *AccountRepository) Create(account *Account) (*Account, error) {
	result := repo.Db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (repo *AccountRepository) Update(account *Account) (*Account, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (repo *AccountRepository) FindByUserId(id uint) (*Account, error) {
	var account Account
	result := repo.Db.First(&account, "user_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}
