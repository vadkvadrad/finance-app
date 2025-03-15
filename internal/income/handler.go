package income

import (
	"finance-app/configs"
	"finance-app/pkg/er"
	"finance-app/pkg/middleware"
	"finance-app/pkg/req"
	"finance-app/pkg/res"
	"net/http"
)

type IncomeHandler struct {
	Config        *configs.Config
	IncomeService *IncomeService
}

type IncomeHandlerDeps struct {
	Config        *configs.Config
	IncomeService *IncomeService
}

func NewIncomeHandler(router *http.ServeMux, deps IncomeHandlerDeps) {
	handler := &IncomeHandler{
		Config:        deps.Config,
		IncomeService: deps.IncomeService,
	}

	router.Handle("POST /income/create", middleware.IsAuthed(handler.Create(), deps.Config))
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
