package income

import (
	"finance-app/configs"
	"finance-app/internal/account"
	"finance-app/pkg/er"
	"finance-app/pkg/middleware"
	"finance-app/pkg/req"
	"finance-app/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type IncomeHandler struct {
	Config         *configs.Config
	IncomeService  *IncomeService
	AccountService *account.AccountService
}

type IncomeHandlerDeps struct {
	Config         *configs.Config
	IncomeService  *IncomeService
	AccountService *account.AccountService
}

func NewIncomeHandler(router *http.ServeMux, deps IncomeHandlerDeps) {
	handler := &IncomeHandler{
		Config:         deps.Config,
		IncomeService:  deps.IncomeService,
		AccountService: deps.AccountService,
	}

	router.Handle("POST /income", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("DELETE /income/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.Handle("PATCH /income/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
}

func (handler *IncomeHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, ok := r.Context().Value(middleware.ContextUserDataKey).(middleware.UserData)
		if !ok {
			http.Error(w, er.ErrWrongUserCredentials, http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[NewIncomeRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		income, err := handler.IncomeService.NewIncome(&Income{
			UserId: userData.Id,
			Amount: body.Amount,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, income, http.StatusCreated)
	}
}

func (handler *IncomeHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, ok := r.Context().Value(middleware.ContextUserDataKey).(middleware.UserData)
		if !ok {
			http.Error(w, er.ErrWrongUserCredentials, http.StatusUnauthorized)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.IncomeService.DeleteIncome(uint(id), userData.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (handler *IncomeHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение данных юзера
		userData, ok := r.Context().Value(middleware.ContextUserDataKey).(middleware.UserData)
		if !ok {
			http.Error(w, er.ErrWrongUserCredentials, http.StatusUnauthorized)
			return
		}

		// Парсинг id у дохода
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Получение обновленного дохода
		body, err := req.HandleBody[NewIncomeRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Обновление старого дохода
		income, err := handler.IncomeService.RedactIncome(&Income{
			Model:  gorm.Model{ID: uint(id)}, // Индекс дохода
			UserId: userData.Id,              // id текущего пользователя
			Amount: body.Amount,              // обновленная цена
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, income, http.StatusOK)
	}
}
