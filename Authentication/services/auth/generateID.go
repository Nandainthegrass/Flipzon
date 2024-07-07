package auth

import (
	"github.com/google/uuid"
)

func GenerateID() string {
	id := uuid.New()
	return id.String()
}
