package account

import (
	"finance-app/configs"
	"finance-app/pkg/er"
	"finance-app/pkg/middleware"
	"finance-app/pkg/res"
	"net/http"
)

type AccountHandler struct {
	Config *configs.Config
	AccountRepository *AccountRepository
}

type AccountHandlerDeps struct {
	Config *configs.Config
	AccountRepository *AccountRepository
}

func NewAccountHandler(router *http.ServeMux, deps AccountHandlerDeps) {
	handler := &AccountHandler{
		Config: deps.Config,
		AccountRepository: deps.AccountRepository,
	}

	router.Handle("GET /account", middleware.IsAuthed(handler.Get(), deps.Config))
}


func (handler *AccountHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение email из контекста
		userData, ok := r.Context().Value(middleware.ContextUserDataKey).(middleware.UserData)
		if !ok {
			http.Error(w, er.ErrWrongUserCredentials, http.StatusBadRequest)
			return
		}

		// Получение данных аккаунта
		acc, err := handler.AccountRepository.FindByUserId(userData.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return 
		}

		res.Json(w, acc, http.StatusOK)
	}
}