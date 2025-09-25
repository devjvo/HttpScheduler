package Controller

import (
	entity "HttpScheduler/src/Api/Domain/Entity"
	factory "HttpScheduler/src/Api/Infrastructure/Factory"
	repository "HttpScheduler/src/Api/Infrastructure/Repository"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ResponseRequestList struct {
	Items      []entity.Request `json:"items"`
	NextCursor *string          `json:"nextCursor"`
}

type RequestController struct {
	uuidFactory factory.UuidFactory
}

func NewRequestController() *RequestController {
	return &RequestController{
		uuidFactory: *factory.NewUuidFactory(),
	}
}

func (r *RequestController) GetRequest(response http.ResponseWriter, request *http.Request) error {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")

	cursorRaw := mux.Vars(request)["id"]
	cursor, err := r.uuidFactory.CreateFromString(cursorRaw)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": fmt.Sprintf("cursor is an invalid uuid: %s", err),
		})

		return nil
	}

	scheduledRequest := repository.NewRequestRepository().Get(cursor)

	if scheduledRequest == nil {
		response.WriteHeader(http.StatusNotFound)
		return nil
	}

	if err := json.NewEncoder(response).Encode(scheduledRequest); err != nil {
		message := fmt.Sprintf("unable to encode request list to json. error: %s", err.Error())
		slog.Error(message)
		panic(message)
	}

	return nil
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
		parsedCursor, err := r.uuidFactory.CreateFromString(cursorRaw)

		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{
				"error": fmt.Sprintf("cursor is an invalid uuid. %s", err),
			})

			return nil
		}

		cursor = parsedCursor
	}

	limit, err := strconv.ParseInt(request.Form.Get("limit"), 10, 8)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": fmt.Sprintf("limit is an invalid integer. %s", request.Form.Get("limit")),
		})

		return nil
	}

	requestList, err := repository.NewRequestRepository().ListRequest(cursor, limit+1)

	if err != nil {
		panic(err)
	}

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
