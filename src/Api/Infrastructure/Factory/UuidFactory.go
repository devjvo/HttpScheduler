package factory

import (
	"fmt"

	"github.com/google/uuid"
)

type UuidFactory struct{}

func NewUuidFactory() *UuidFactory {
	return &UuidFactory{}
}

func (u *UuidFactory) CreateFromString(s string) (uuid.UUID, error) {
	uuidv7, err := uuid.Parse(s)

	if err != nil {
		return uuidv7, fmt.Errorf("uuid %s is malformed: %s", s, err)
	}

	if uuidv7.Version() != 7 {
		return uuidv7, fmt.Errorf("uuid %s must be version: 7, got version: %d", s, uuidv7.Version())
	}

	return uuidv7, nil
}
