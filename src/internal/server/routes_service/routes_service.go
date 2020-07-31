package routes_service

import (
	"github.com/gorilla/mux"
	"github.com/palantir/stacktrace"
	"internal/server/handlers"
	"log"
	"net/http"
)

func InitRoutes() error {
	log.Println("Configuring routes")
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1/").Subrouter()
	subRouter.HandleFunc("/create", handlers.ProcessCreateHandler)
	subRouter.HandleFunc("/receive", handlers.ProcessReceiveHandler)
	http.Handle("/", router)

	log.Println("Server is listening in the port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to init the routes")
	}
	return nil
}
