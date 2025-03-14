package req

import (
	"finance-app/pkg/res"
	"net/http"
)

func HandleBody[T any] (w http.ResponseWriter, r *http.Request) (*T, error){
	body, err := Decode[T](r.Body)
	if err != nil {
		res.Json(w, err.Error(), http.StatusPaymentRequired)
		return nil, err
	}
	
	err = IsValid(body)
	if err != nil {
		res.Json(w, err.Error(), http.StatusPaymentRequired)
		return nil, err
	}
	return &body, nil
}
