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

type CreateTenderRequest struct {
	Name           string
	Description    string
	ServiceType    string
	OrganizationID uuid.UUID
	Creator        string
}

type CreateBidsRequest struct {
	Name        string
	Description string
	TenderID    uuid.UUID
	AuthorType  string
	AuthorID    uuid.UUID
}

type CreateBidsResponse struct {
	ID          uuid.UUID
	Name        string
	Status      string
	Description string
	AuthorType  string
	AuthorID    uuid.UUID
	Version     int
	CreatedAt   time.Time
}

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EditTenderParams struct {
	TenderID    uuid.UUID
	Username    string
	Name        string
	Description string
	ServiceType string
}
