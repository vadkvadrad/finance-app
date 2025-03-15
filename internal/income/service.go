package income

import (
	"errors"
	"finance-app/internal/account"
	"finance-app/pkg/er"
)

type IncomeService struct {
	IncomeRepository  *IncomeRepository
	AccountRepository *account.AccountRepository
}

type IncomeServiceDeps struct {
	IncomeRepository *IncomeRepository
	AccountRepository *account.AccountRepository
}


func NewIncomeService(deps IncomeServiceDeps) *IncomeService {
	return &IncomeService{
		IncomeRepository: deps.IncomeRepository,
		AccountRepository: deps.AccountRepository,
	}
}


func (service *IncomeService) NewIncome(income *Income) (*Income, error) {
	// Проверка на положительность дохода
	if income.Amount <= 0 {
		return nil, errors.New(er.ErrNegativeIncome)
	}

	// Создание нового дохода
	createdIncome, err := service.IncomeRepository.Create(income)
	if err != nil {
		return nil, err
	}

	// Получение аккаунта
	account, err := service.AccountRepository.FindByUserId(createdIncome.UserId)
	if err != nil {
		return nil, err
	}

	// Обновление баланса аккаунта
	account.Balance += income.Amount
	_, err = service.AccountRepository.Update(account)
	if err != nil {
		return nil, err
	}

	return createdIncome, nil
}

func (service *IncomeService) RedactIncome(income *Income) (*Income, error) {
	return service.IncomeRepository.Update(income)
}


