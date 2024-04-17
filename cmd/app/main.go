package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/lincolngondin/desafio-criptografia/config"
	"github.com/lincolngondin/desafio-criptografia/internal/user"
	_ "modernc.org/sqlite"
)

func main() {
	configs := config.New()

	db, err := sql.Open(configs.DBDriverName, configs.DBDataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUserHandler)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserHandler)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserHandler)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUserHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
