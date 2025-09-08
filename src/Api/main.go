package main

import (
	"HttpScheduler/src/Api/Ui/Http/Controller"
	security "HttpScheduler/src/Api/Ui/Http/Security"
	"log"
	"net/http"
)

func main() {
	basicAuthChecker := security.NewBasicAuthChecker()

	http.HandleFunc(
		"/v1/healthcheck",
		basicAuthChecker.Middleware(
			func(response http.ResponseWriter, request *http.Request) {
				if request.Method == http.MethodHead {
					Controller.HealthCheck(response, request)
					return
				}

				response.WriteHeader(http.StatusMethodNotAllowed)
			}))

	http.HandleFunc(
		"/v1/request",
		basicAuthChecker.Middleware(
			func(response http.ResponseWriter, request *http.Request) {
				if request.Method == http.MethodGet {
					Controller.NewRequestController().ListRequest(response, request)
					return
				}

				response.WriteHeader(http.StatusMethodNotAllowed)
			}))

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
