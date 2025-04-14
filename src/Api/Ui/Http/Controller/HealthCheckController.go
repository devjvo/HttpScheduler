package Controller

import "net/http"

func HealthCheck(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodHead {
		response.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
}
