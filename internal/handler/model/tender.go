package model

import (
	"github.com/google/uuid"
	"time"
)

type CreateRequest struct {
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description" validate:"required"`
	ServiceType    string    `json:"service_type" validate:"oneof=Delivery Manufacture Construction"`
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	Creator        string    `json:"creator_username" validate:"required"`
}

type StatusRequest struct {
	TenderID uuid.UUID `validate:"required"`
	Username string
}

type EditStatusRequest struct {
	TenderID uuid.UUID `validate:"required"`
	Username string    `validate:"required"`
	Status   string    `validate:"oneof=Created Published Closed"`
}

type ListFilter struct {
	Limit       string
	Offset      string
	ServiceType []string `validate:"oneof=Delivery Manufacture Construction"`
}

type UserListFilter struct {
	Limit    string
	Offset   string
	Username string `validate:"required"`
}

type IntUserListFilter struct {
	Limit    int32
	Offset   int32
	Username string `validate:"required"`
}

type IntListFilter struct {
	Limit       int32
	Offset      int32
	ServiceType []string
}

type Tender struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      string
	ServiceType string
	Version     int
	CreatedAt   time.Time
}
