package app

import "net/http"

func App() http.Handler {
	router := http.NewServeMux()


	return router
}