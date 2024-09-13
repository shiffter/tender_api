package model

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	Creator        string    `db:"creator_username"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Status         string    `db:"status"`
	ServiceType    string    `db:"service_type"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type EditTenderParams struct {
	Name        string
	Description string
	ServiceType string
}
