package main

import (
	"net/http"
	"time"

	"github.com/WilkerAlves/go-expert-chalanger-client-server-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Second))

	handlerCep := handlers.NewCepHandler()
	router.Get("/{cep}", handlerCep.GetCep)

	http.ListenAndServe(":8080", router)
}
