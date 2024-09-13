package converter

import (
	"github.com/google/uuid"
	"strconv"
	apiModel "tender/internal/handler/model"
	serviceModel "tender/internal/service/model"
)

func FromStringToIntListFilter(req *apiModel.ListFilter) (*apiModel.IntListFilter, error) {
	if req.Offset == "" {
		req.Offset = "0"
	}
	o64, err := strconv.Atoi(req.Offset)
	if err != nil {
		return nil, err
	}

	if req.Limit == "" {
		req.Limit = "5"
	}
	l64, err := strconv.Atoi(req.Limit)
	if err != nil {
		return nil, err
	}

	return &apiModel.IntListFilter{
		Limit:       int32(l64),
		Offset:      int32(o64),
		ServiceType: req.ServiceType,
	}, nil
}

func FromStringToIntUserListFilter(req *apiModel.UserListFilter) (*apiModel.IntUserListFilter, error) {
	if req.Offset == "" {
		req.Offset = "0"
	}
	o64, err := strconv.Atoi(req.Offset)
	if err != nil {
		return nil, err
	}

	if req.Limit == "" {
		req.Limit = "5"
	}
	l64, err := strconv.Atoi(req.Limit)
	if err != nil {
		return nil, err
	}

	return &apiModel.IntUserListFilter{
		Limit:    int32(l64),
		Offset:   int32(o64),
		Username: req.Username,
	}, nil
}

func FromApiToServiceEditTenderParam(api *apiModel.EditTenderRequest) (*serviceModel.EditTenderParams, error) {
	tenderUID, err := uuid.Parse(api.TenderID)
	if err != nil {
		return nil, err
	}

	return &serviceModel.EditTenderParams{
		TenderID:    tenderUID,
		Username:    api.Username,
		Name:        api.Params.Name,
		Description: api.Params.Description,
		ServiceType: api.Params.ServiceType,
	}, nil
}
