package main

import (
	"finance-app/cmd/app"
	"fmt"
	"net/http"
)

func main() {
	app := app.App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("server started on port", server.Addr)
	server.ListenAndServe()
}