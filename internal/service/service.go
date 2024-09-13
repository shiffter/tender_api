package service

import (
	"context"
	"github.com/google/uuid"
	apiModel "tender/internal/handler/model"
	serviceModel "tender/internal/service/model"
)

type TenderService interface {
	Create(ctx context.Context, request *apiModel.CreateRequest) (*apiModel.Tender, error)
	Status(ctx context.Context, id uuid.UUID, username string) (string, error)
	EditStatus(ctx context.Context, tenderID uuid.UUID, username, status string) (*apiModel.Tender, error)
	EditTender(ctx context.Context, param *apiModel.EditTenderRequest) (*apiModel.Tender, error)
	List(ctx context.Context, limit, offset int32, serviceTypes []string) ([]*apiModel.Tender, error)
	UserList(ctx context.Context, limit, offset int32, username string) ([]*apiModel.Tender, error)
}

type UserService interface {
	Get(ctx context.Context, username string) (*serviceModel.User, error)
}

type BidsService interface {
	Create(ctx context.Context, req *apiModel.CreateBidRequest) (*apiModel.CreateBidResp, error)
}
