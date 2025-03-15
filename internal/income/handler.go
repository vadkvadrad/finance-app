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
)

type IncomeHandler struct {
	Config        *configs.Config
	IncomeService *IncomeService
	AccountService *account.AccountService
}

type IncomeHandlerDeps struct {
	Config        *configs.Config
	IncomeService *IncomeService
	AccountService *account.AccountService
}

func NewIncomeHandler(router *http.ServeMux, deps IncomeHandlerDeps) {
	handler := &IncomeHandler{
		Config:        deps.Config,
		IncomeService: deps.IncomeService,
		AccountService: deps.AccountService,
	}

	router.Handle("POST /income", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("DELETE /income/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
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