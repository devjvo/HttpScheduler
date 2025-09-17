package main

import (
	"HttpScheduler/src/Api/Ui/Http/Controller"
	security "HttpScheduler/src/Api/Ui/Http/Security"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	basicAuthChecker := security.NewBasicAuthChecker()
	router := mux.NewRouter()

	router.HandleFunc(
		"/v1/healthcheck",
		basicAuthChecker.Middleware(
			func(response http.ResponseWriter, request *http.Request) {
				if request.Method == http.MethodHead {
					Controller.HealthCheck(response, request)
					return
				}

				response.WriteHeader(http.StatusMethodNotAllowed)
			}))

	router.HandleFunc(
		"/v1/request",
		basicAuthChecker.Middleware(
			func(response http.ResponseWriter, request *http.Request) {
				if request.Method == http.MethodGet {
					Controller.NewRequestController().ListRequest(response, request)
					return
				}

				response.WriteHeader(http.StatusMethodNotAllowed)
			}))

	router.HandleFunc(
		"/v1/request/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}",
		basicAuthChecker.Middleware(
			func(response http.ResponseWriter, request *http.Request) {
				if request.Method == http.MethodGet {
					Controller.NewRequestController().GetRequest(response, request)
					return
				}

				response.WriteHeader(http.StatusMethodNotAllowed)
			}))

	http.ListenAndServe(":8080", router)
}
