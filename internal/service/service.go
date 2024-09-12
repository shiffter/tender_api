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
	List(ctx context.Context, limit, offset int32, serviceTypes []string) ([]*apiModel.Tender, error)
	UserList(ctx context.Context, limit, offset int32, username string) ([]*apiModel.Tender, error)
}

type UserService interface {
	Get(ctx context.Context, username string) (*serviceModel.User, error)
}

//type Service struct {
//	repo storage.Storage
//}
//
//func NewService(storage storage.Storage) *Service {
//	return &Service{
//		Storage: storage,
//	}
//}
//
//func (s *Service) Create(ctx context.Context, params CreateTenderRequest) (*CreateTenderResponse, error) {
//	var (
//		storageParams = storage.CreateRequest{
//			Name:           params.Name,
//			Description:    params.Description,
//			ServiceType:    params.ServiceType,
//			OrganizationID: params.OrganizationID,
//			Creator:        params.Creator,
//		}
//	)
//
//	tender, err := s.Storage.Create(ctx, storageParams)
//	if err != nil {
//		return nil, err
//	}
//
//	return &CreateTenderResponse{
//		ID:          tender.ID,
//		Name:        tender.Name,
//		Description: tender.Description,
//		Status:      tender.Status,
//		ServiceType: tender.ServiceType,
//		Version:     tender.Version,
//		CreatedAt:   tender.CreatedAt,
//	}, nil
//}
