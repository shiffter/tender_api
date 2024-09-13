package bids

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	apiModel "tender/internal/handler/model"
	"tender/internal/repository"
	"tender/internal/service"
	serviceModel "tender/internal/service/model"
	"tender/internal/service/model/converter"
	"tender/pkg/stderrs"
)

type bidsService struct {
	bidsRepo   repository.BidsRepos
	userRepo   repository.UsersRepos
	tenderRepo repository.TenderRepos
}

func NewBidsService(bids repository.BidsRepos,
	user repository.UsersRepos,
	tRepo repository.TenderRepos) service.BidsService {
	return &bidsService{
		bidsRepo:   bids,
		userRepo:   user,
		tenderRepo: tRepo,
	}
}

func (s *bidsService) Create(ctx context.Context, req *apiModel.CreateBidRequest) (*apiModel.CreateBidResp, error) {

	userUUID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrUserNotFound
		}

		return nil, err
	}

	tenderUUID, err := uuid.Parse(req.TenderID)
	if err != nil {
		return nil, err
	}

	tender, err := s.tenderRepo.Get(ctx, tenderUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrTenderNotFound
		}

		return nil, err
	}

	if tender.Status == "CREATED" || tender.Status == "CLOSED" {
		_, err = s.tenderRepo.OrganizationRightsForUser(ctx, user.Username, tender.OrganizationId)
		if err != nil {
			return nil, err
		}
	}

	bid, err := s.bidsRepo.Create(ctx, &serviceModel.CreateBidsRequest{
		Name:        req.Name,
		Description: req.Description,
		TenderID:    tenderUUID,
		AuthorType:  req.AuthorType,
		AuthorID:    userUUID,
	})
	if err != nil {
		return nil, err
	}

	return converter.FromServiceToApiBidsResp(bid), nil
}
