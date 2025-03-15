package income

import (
	"finance-app/pkg/db"

	"gorm.io/gorm/clause"
)

type IncomeRepository struct {
	Db *db.Db
}

func NewIncomeRepository(db *db.Db) *IncomeRepository {
	return &IncomeRepository{
		Db: db,
	}
}

func (repo *IncomeRepository) Create(income *Income) (*Income, error) {
	result := repo.Db.Create(income)
	if result.Error != nil {
		return nil, result.Error
	}
	return income, nil
}

func (repo *IncomeRepository) Update(income *Income) (*Income, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(income)
	if result.Error != nil {
		return nil, result.Error
	}
	return income, nil
}

func (repo *IncomeRepository) Delete(id uint) error {
	result := repo.Db.Delete(&Income{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *IncomeRepository) FindByUserId(id uint) (*Income, error) {
	var income Income
	result := repo.Db.First(&income,"user_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &income, nil
}

func (repo *IncomeRepository) FindById(id uint) (*Income, error) {
	var income Income
	result := repo.Db.First(&income,"id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &income, nil
}