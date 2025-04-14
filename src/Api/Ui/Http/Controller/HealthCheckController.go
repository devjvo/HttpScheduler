package Controller

import "net/http"

func HealthCheck(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNoContent)
}
