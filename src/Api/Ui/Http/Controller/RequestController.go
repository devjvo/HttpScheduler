package Controller

import (
	"fmt"
	"net/http"
)

func ListRequest(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)

	fmt.Fprintf(response, "Request list")
}
