package repository

import (
	entity "HttpScheduler/src/Api/Domain/Entity"
	database "HttpScheduler/src/Api/Infrastructure/Database"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository() *RequestRepository {
	return &RequestRepository{
		db: database.GetDatabaseConnection(),
	}
}

func (r *RequestRepository) Get(id uuid.UUID) *entity.Request {
	row := r.db.QueryRow(
		`SELECT id, created_at, http_method, url, response_code FROM request WHERE id = $1`,
		id,
	)

	var request entity.Request

	if err := row.Scan(&request.Id, &request.CreatedAt, &request.HttpMethod, &request.Url, &request.ResponseCode); err != nil {
		if err != sql.ErrNoRows {
			slog.Error(fmt.Sprintf("unable to get row in database for request id: %s. error: %s", id.String(), err.Error()))
		}

		return nil
	}

	return &request
}

func (r *RequestRepository) ListRequest(cursor uuid.UUID, limit int64) []entity.Request {
	rows, err := r.db.Query(
		`SELECT id, created_at, http_method, url, response_code
		FROM request
		WHERE id >= $1
		ORDER BY id
		LIMIT $2`,
		cursor,
		limit,
	)

	if err != nil {
		message := fmt.Sprintf("unable to select request entity. error: %s", err.Error())
		slog.Error(message)
		panic(message)
	}

	var requestList []entity.Request

	for rows.Next() {
		var request entity.Request

		if err := rows.Scan(&request.Id, &request.CreatedAt, &request.HttpMethod, &request.Url, &request.ResponseCode); err != nil {
			message := fmt.Sprintf("unable to scan request entity. error: %s", err.Error())
			slog.Error(message)
			panic(message)
		}

		requestList = append(requestList, request)
	}

	return requestList
}
