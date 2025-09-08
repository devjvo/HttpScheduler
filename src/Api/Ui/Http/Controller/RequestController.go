package Controller

import (
	entity "HttpScheduler/src/Api/Domain/Entity"
	repository "HttpScheduler/src/Api/Infrastructure/Repository"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type ResponseRequestList struct {
	Items      []entity.Request `json:"items"`
	NextCursor *string          `json:"nextCursor"`
}

type RequestController struct{}

func NewRequestController() *RequestController {
	return &RequestController{}
}

func (r *RequestController) ListRequest(response http.ResponseWriter, request *http.Request) error {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")

	defer func() {
		if recover := recover(); recover != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(map[string]string{
				"error": "internal server error",
			})
		}
	}()

	request.ParseForm()

	cursorRaw := request.Form.Get("cursor")
	cursor := uuid.Nil

	if cursorRaw != "" {
		parsedCursor, err := uuid.Parse(cursorRaw)

		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{
				"error": fmt.Sprintf("cursor is an invalid uuid: %s", cursorRaw),
			})

			return nil
		}

		if parsedCursor.Version() != 7 {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{
				"error": fmt.Sprintf("cursor is uuid but must be v7: %s", cursorRaw),
			})

			return nil
		}

		cursor = parsedCursor
	}

	limit, err := strconv.ParseInt(request.Form.Get("limit"), 10, 8)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": fmt.Sprintf("limit is an invalid integer: %s", request.Form.Get("limit")),
		})

		return nil
	}

	requestList := repository.NewRequestRepository().ListRequest(cursor, limit+1)
	var nextCursor *string

	if len(requestList) > int(limit) {
		lastId := requestList[len(requestList)-1].Id
		requestList = requestList[:limit]
		nextCursor = &lastId
	}

	responseBody := ResponseRequestList{requestList, nextCursor}

	if err := json.NewEncoder(response).Encode(responseBody); err != nil {
		message := fmt.Sprintf("unable to encode request list to json. error: %s", err.Error())
		slog.Error(message)
		panic(message)
	}

	return nil
}
