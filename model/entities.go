package model

import (
	"github.com/google/uuid"
	"time"
)

type Annotation struct {
	ID          uuid.UUID
	VideoURL    string
	Start       time.Duration
	End         time.Duration
	UserCreated uuid.UUID
	Name        string
	Comment     string
}

type Video struct {
	URL         string
	Duration    time.Duration
	UserCreated uuid.UUID
}

type User struct {
	Name           string
	ID             uuid.UUID
	HashedPassword string
}
