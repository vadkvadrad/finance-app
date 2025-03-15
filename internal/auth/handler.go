package auth

import (
	"finance-app/configs"
	"finance-app/pkg/jwt"
	"finance-app/pkg/req"
	"finance-app/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService *AuthService
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService *AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/verify", handler.Verify())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			return
		}

		user, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJwt(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Id:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// отправить ответ
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		sessionId, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// отправить ответ
		data := RegisterResponse{
			SessionId: sessionId,
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := handler.AuthService.Verify(body.SessionId, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := jwt.NewJwt(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Id:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := VerifyResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusAccepted)
	}
}
