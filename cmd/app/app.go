package app

import (
	"finance-app/configs"
	"finance-app/internal/auth"
	"finance-app/internal/user"
	"finance-app/pkg/db"
	"finance-app/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf, err := configs.Load()
	if err != nil {
		panic(err)
	}
	db := db.NewDb(conf)
	router := http.NewServeMux()
	//eventBus := event.NewEventBus()

	// Repository
	userRepository := user.NewUserRepository(db)

	// Services 
	authService := auth.NewAuthService(userRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
		AuthService: authService,
	})


	// listening for statistic
	//go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.Logging,
	)

	return stack(router)
}