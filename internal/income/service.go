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
	IncomeRepository  *IncomeRepository
	AccountRepository *account.AccountRepository
}

func NewIncomeService(deps IncomeServiceDeps) *IncomeService {
	return &IncomeService{
		IncomeRepository:  deps.IncomeRepository,
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
	// Проверка на положительность баланса
	if income.Amount <= 0 {
		return nil, errors.New(er.ErrNegativeIncome)
	}

	// Получение старого дохода
	oldIncome, err := service.IncomeRepository.FindById(income.ID)
	if err != nil {
		return nil, err
	}

	// Проверка на подлинность пользователя
	if oldIncome.UserId != income.UserId {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}

	// Получение аккаунта с данным доходом
	account, err := service.AccountRepository.FindByUserId(oldIncome.UserId)
	if err != nil {
		return nil, err
	}

	// Обновление баланса аккаунта
	account.Balance -= oldIncome.Amount
	account.Balance += income.Amount

	// Обновление аккаунта
	_, err = service.AccountRepository.Update(account)
	if err != nil {
		return nil, err
	}

	// Обновление дохода
	income, err = service.IncomeRepository.Update(income)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (service *IncomeService) DeleteIncome(id uint, userId uint) error {
	// Получение дохода
	income, err := service.IncomeRepository.FindById(id)
	if err != nil {
		return err
	}

	// Получение аккаунта
	account, err := service.AccountRepository.FindByUserId(userId)
	if err != nil {
		return err
	}

	// Проверка правильный ли получен аккаунт
	if income.UserId != account.UserID {
		return errors.New(er.ErrWrongUserCredentials)
	}

	// Изменение баланса
	account.Balance -= income.Amount

	// Обновление аккаунта
	_, err = service.AccountRepository.Update(account)
	if err != nil {
		return err
	}

	// Удаление дохода
	err = service.IncomeRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
