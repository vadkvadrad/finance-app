package auth

import (
	"errors"
	"finance-app/internal/account"
	"finance-app/internal/user"
	"finance-app/pkg/er"
	"finance-app/pkg/event"
	"finance-app/pkg/sender"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository    *user.UserRepository
	AccountRepository *account.AccountRepository
	Event             *event.EventBus
}

type AuthServiceDeps struct {
	UserRepository    *user.UserRepository
	AccountRepository *account.AccountRepository
	Event             *event.EventBus
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		UserRepository:    deps.UserRepository,
		AccountRepository: deps.AccountRepository,
		Event:             deps.Event,
	}
}

func (service *AuthService) Login(email, password string) (*user.User, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(user.EmailKey, email)
	if existedUser == nil {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}

	// Проверка, верифицирован ли пользователь
	if !existedUser.IsVerified {
		return nil, errors.New(er.ErrUserNotVerified)
	}

	// Сравниваем пароли
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}

	return existedUser, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(user.EmailKey, email)

	if existedUser != nil && existedUser.IsVerified { // если пользователь существует и верифицирован
		return "", errors.New(er.ErrUserExists)
	} else if existedUser != nil { // если пользователь существует и НЕ верифицирован
		// Регенерация кода и id сессии
		existedUser.Generate()

		// Отправить сообщение с кодом
		service.sendEmail(email, existedUser.Code)

		// Перезаписать юзера в бд
		_, err := service.UserRepository.Update(existedUser)
		if err != nil {
			return "", err
		}

		return existedUser.SessionId, nil
	}

	// Генерим хеш пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Создаем модель юзера
	user := &user.User{
		Email:      email,
		Password:   string(hashedPassword),
		Name:       name,
		Role:       user.RoleUser,
		IsVerified: false,
	}
	user.Generate()

	// Создаем запись user
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}

	// Создаем запись account
	_, err = service.AccountRepository.Create(&account.Account{
		UserID:   user.ID,
		Balance:  0,
		Currency: account.CurrencyRub,
	})
	if err != nil {
		return "", err
	}

	// Отправить сообщение с кодом
	service.sendEmail(email, user.Code)

	return user.SessionId, nil
}

func (service *AuthService) Verify(sessionId, code string) (*user.User, error) {
	// Находим пользователя
	existedUser, _ := service.UserRepository.FindByKey(user.SessionIdKey, sessionId)
	if existedUser == nil {
		return nil, errors.New(er.ErrUserExists)
	}

	// Проверка на подлинность кода
	if existedUser.Code != code {
		return nil, errors.New(er.ErrNotAuthorized)
	}

	// Пользователь становится верифицированным
	existedUser.IsVerified = true
	user, err := service.UserRepository.Update(existedUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *AuthService) sendEmail(email, code string) {
	// Отправка кода на почту
	go service.Event.Publish(event.Event{
		Type: event.EventSendEmail,
		Data: sender.Addressee{
			To:      email,
			Subject: "Подтвердите почту",
			Text:    "Ваш персональный код подтверждения личности: " + code + ". Не сообщайте никому данный код.",
		},
	})
}
