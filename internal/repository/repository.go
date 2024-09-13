package repository

import (
	"context"
	"github.com/google/uuid"
	serviceModel "tender/internal/service/model"
)

type TenderRepos interface {
	Create(ctx context.Context, tender *serviceModel.CreateTenderRequest) (*serviceModel.Tender, error)
	List(ctx context.Context, limit, offset int32, serviceTypes []string) ([]*serviceModel.Tender, error)
	UserList(ctx context.Context, limit, offset int32, serviceTypes string) ([]*serviceModel.Tender, error)
	OrganizationRightsForUser(ctx context.Context, userName string, orgUUID uuid.UUID) (*uuid.UUID, error)
	Get(ctx context.Context, tenderUUID uuid.UUID) (*serviceModel.Tender, error)
	EditStatus(ctx context.Context, status string, tenderID uuid.UUID) (*serviceModel.Tender, error)
	EditTender(c context.Context, p *serviceModel.EditTenderParams, tenderID uuid.UUID) (*serviceModel.Tender, error)
}

type BidsRepos interface {
	Create(ctx context.Context, bids *serviceModel.CreateBidsRequest) (*serviceModel.CreateBidsResponse, error)
}

type UsersRepos interface {
	Get(ctx context.Context, username string) (*serviceModel.User, error)
	GetByID(ctx context.Context, userUUID uuid.UUID) (*serviceModel.User, error)
}
