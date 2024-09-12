package converter

import (
	"strconv"
	"tender/internal/handler/model"
)

func FromStringToIntListFilter(req *model.ListFilter) (*model.IntListFilter, error) {
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

	return &model.IntListFilter{
		Limit:       int32(l64),
		Offset:      int32(o64),
		ServiceType: req.ServiceType,
	}, nil

}

func FromStringToIntUserListFilter(req *model.UserListFilter) (*model.IntUserListFilter, error) {
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

	return &model.IntUserListFilter{
		Limit:    int32(l64),
		Offset:   int32(o64),
		Username: req.Username,
	}, nil

}
