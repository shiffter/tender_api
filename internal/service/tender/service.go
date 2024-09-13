package tender

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
	apiModel "tender/internal/handler/model"
	"tender/internal/repository"
	svc "tender/internal/service"
	"tender/internal/service/model/converter"
	"tender/pkg/stderrs"
)

type tenderService struct {
	tenderRepo repository.TenderRepos
	userRepo   repository.UsersRepos
}

func NewTenderService(tenderRepo repository.TenderRepos, userRepo repository.UsersRepos) svc.TenderService {
	return &tenderService{
		tenderRepo: tenderRepo,
		userRepo:   userRepo,
	}
}

func (s *tenderService) Create(ctx context.Context, req *apiModel.CreateRequest) (*apiModel.Tender, error) {

	if _, err := s.tenderRepo.OrganizationRightsForUser(ctx, req.Creator, req.OrganizationID); err != nil {
		return nil, err
	}

	svcModel := converter.FromApiToServiceCreateReq(req)

	svcTender, err := s.tenderRepo.Create(ctx, svcModel)
	if err != nil {
		return nil, err
	}

	return converter.FromServiceToApi(svcTender), nil
}

func (s *tenderService) List(ctx context.Context, limit, offset int32, serviceTypes []string) ([]*apiModel.Tender, error) {
	tenders, err := s.tenderRepo.List(ctx, limit, offset, serviceTypes)
	if err != nil {
		return nil, err
	}

	tendersForApi := make([]*apiModel.Tender, 0, len(tenders))

	for _, t := range tenders {
		tendersForApi = append(tendersForApi, converter.FromServiceToApi(t))
	}

	return tendersForApi, nil
}

func (s *tenderService) UserList(ctx context.Context, limit, offset int32, username string) ([]*apiModel.Tender, error) {

	_, err := s.userRepo.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	tenders, err := s.tenderRepo.UserList(ctx, limit, offset, username)
	if err != nil {
		return nil, err
	}

	tendersForApi := make([]*apiModel.Tender, 0, len(tenders))

	for _, t := range tenders {
		tendersForApi = append(tendersForApi, converter.FromServiceToApi(t))
	}

	return tendersForApi, nil
}

func (s *tenderService) Status(ctx context.Context, tenderID uuid.UUID, username string) (string, error) {

	user, err := s.userRepo.Get(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", stderrs.ErrUserNotFound
		}

		return "", err
	}

	tender, err := s.tenderRepo.Get(ctx, tenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", stderrs.ErrTenderNotFound
		}

		return "", err
	}

	if user.Username != tender.Creator {
		return "", stderrs.ErrNotEnoughRights
	}

	return tender.Status, nil
}

func (s *tenderService) EditStatus(
	ctx context.Context,
	tenderID uuid.UUID,
	username, status string) (*apiModel.Tender, error) {

	user, err := s.userRepo.Get(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrUserNotFound
		}

		return nil, err
	}

	tender, err := s.tenderRepo.Get(ctx, tenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrTenderNotFound
		}

		return nil, err
	}

	if user.Username != tender.Creator {
		return nil, stderrs.ErrNotEnoughRights
	}

	if tender.Status == strings.ToUpper(status) {
		return converter.FromServiceToApi(tender), nil
	}

	updTender, err := s.tenderRepo.EditStatus(ctx, strings.ToUpper(status), tenderID)
	if err != nil {
		return nil, err
	}

	return converter.FromServiceToApi(updTender), nil
}

func (s *tenderService) EditTender(
	ctx context.Context,
	request *apiModel.EditTenderRequest) (*apiModel.Tender, error) {

	svcParams, err := converter.FromApiToServiceEditTender(request)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Get(ctx, svcParams.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrUserNotFound
		}

		return nil, err
	}

	tender, err := s.tenderRepo.Get(ctx, svcParams.TenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, stderrs.ErrTenderNotFound
		}

		return nil, err
	}

	if user.Username != tender.Creator {
		return nil, stderrs.ErrNotEnoughRights
	}

	// TODO need validate for unchanged params in tender, bug: useless increment version

	updTender, err := s.tenderRepo.EditTender(ctx, svcParams, svcParams.TenderID)
	if err != nil {
		return nil, err
	}

	return converter.FromServiceToApi(updTender), nil
}
