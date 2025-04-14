package main

import (
	"HttpScheduler/src/Api/Ui/Http/Controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(
		"/v1/healthcheck",
		func(response http.ResponseWriter, request *http.Request) {
			Controller.HealthCheck(response, request)
		})

	http.HandleFunc(
		"/v1/request",
		func(response http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				Controller.ListRequest(response, request)
			}
		})

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
