package auth

import (
	"errors"
	"finance-app/internal/user"
	"finance-app/pkg/er"
	"finance-app/pkg/event"
	"finance-app/pkg/sender"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
	Event          *event.EventBus
}

type AuthServiceDeps struct {
	UserRepository *user.UserRepository
	Event          *event.EventBus
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		UserRepository: deps.UserRepository,
		Event:          deps.Event,
	}
}

func (service *AuthService) Login(email, password string) (string, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(user.EmailKey, email)
	if existedUser == nil {
		return "", errors.New(er.ErrWrongUserCredentials)
	}

	// Проверка, верифицирован ли пользователь
	if !existedUser.IsVerified {
		return "", errors.New(er.ErrUserNotVerified)
	}

	// Сравниваем пароли
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(er.ErrWrongUserCredentials)
	}

	return email, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(user.EmailKey, email)
	if existedUser != nil {
		return "", errors.New(er.ErrUserExists)
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

	// Создаем запись в базе данных
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}

	// Отправка кода на почту
	go service.Event.Publish(event.Event{
		Type: event.EventSendEmail,
		Data: sender.Addressee{
			To:      email,
			Subject: "Подтвердите почту",
			Text:    "Ваш персональный код подтверждения личности: " + user.Code + ". Не сообщайте никому данный код.",
		},
	})

	return user.SessionId, nil
}

func (service *AuthService) Verify(sessionId, code string) (string, error) {
	// Находим пользователя
	existedUser, _ := service.UserRepository.FindByKey(user.SessionIdKey, sessionId)
	if existedUser == nil {
		return "", errors.New(er.ErrUserExists)
	}

	// Проверка на подлинность кода
	if existedUser.Code != code {
		return "", errors.New(er.ErrNotAuthorized)
	}

	// Пользователь становится верифицированным
	existedUser.IsVerified = true
	user, err := service.UserRepository.Update(existedUser)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}
