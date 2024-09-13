package bids

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	handler "tender/internal/handler/api/http"
	apiModel "tender/internal/handler/model"
	svc "tender/internal/service"
	"tender/pkg/stderrs"
)

type bidsHandler struct {
	service   svc.BidsService
	validator *validator.Validate
}

func NewBidsHandler(service svc.BidsService) handler.BidsHandler {
	return &bidsHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *bidsHandler) CreateBid(w http.ResponseWriter, r *http.Request) {

	var (
		req = apiModel.CreateBidRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.validator.Struct(req); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.Create(context.Background(), &req)

	if err != nil {
		switch {
		case errors.Is(err, stderrs.ErrUserNotFound):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, stderrs.ErrTenderNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, stderrs.ErrNotEnoughRights):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(err.Error()))

		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
