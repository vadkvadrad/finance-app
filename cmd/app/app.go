package app

import (
	"finance-app/configs"
	"finance-app/internal/auth"
	"finance-app/internal/user"
	"finance-app/pkg/db"
	"finance-app/pkg/event"
	"finance-app/pkg/middleware"
	"finance-app/pkg/sender"
	"net/http"
)

func App() http.Handler {
	conf, err := configs.Load()
	if err != nil {
		panic(err)
	}
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()
	sender, err := sender.Load(conf, eventBus)
	if err != nil {
		panic(err)
	}

	// Repository
	userRepository := user.NewUserRepository(db)

	// Services
	authService := auth.NewAuthService(auth.AuthServiceDeps{
		UserRepository: userRepository,
		Event:          eventBus,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	// listening for statistic
	go sender.Listen()

	// Middlewares
	stack := middleware.Chain(
		middleware.Logging,
	)

	return stack(router)
}
