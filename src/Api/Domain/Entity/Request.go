package entity

import (
	"time"
)

type Request struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	HttpMethod   string    `json:"httpMethod"`
	Url          string    `json:"url"`
	ResponseCode *string   `json:"responseCode"`
}
