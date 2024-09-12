package model

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Status         string
	Creator        string
	OrganizationId uuid.UUID
	ServiceType    string
	Version        int
	CreatedAt      time.Time
}

type CreateRequest struct {
	Name           string
	Description    string
	ServiceType    string
	OrganizationID uuid.UUID
	Creator        string
}

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
