package account

type AccountService struct {
	AccountRepository *AccountRepository
}

type AccountServiceDeps struct {
	AccountRepository *AccountRepository
}

func NewAccountService(deps AccountServiceDeps) *AccountService {
	return &AccountService{
		AccountRepository: deps.AccountRepository,
	}
}


func (service *AccountService) GetByUserId(id uint) (*Account, error) {
	return service.AccountRepository.FindByUserId(id)
}