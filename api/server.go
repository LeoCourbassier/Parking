package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Start starts the webserver
func Start() {
	logrus.Info("Starting server...")
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	NewParkingRouter(router)

	recoveryRouter := handlers.RecoveryHandler()(router)

	err := http.ListenAndServe(":4000", recoveryRouter)
	if err != nil {
		panic(err)
	}
}
