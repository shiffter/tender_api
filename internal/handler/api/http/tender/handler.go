package tender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	handler "tender/internal/handler/api/http"
	apiModel "tender/internal/handler/model"
	"tender/internal/handler/model/converter"
	"tender/internal/service"
	"tender/pkg/stderrs"
)

type Handler struct {
	service   service.TenderService
	validator *validator.Validate
}

func NewHandler(service service.TenderService) handler.TenderHandler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var (
		req = apiModel.CreateRequest{}
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	req.ServiceType = strings.ToUpper(req.ServiceType)

	if err := h.validator.Struct(req); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.Create(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	var (
		filter               = apiModel.ListFilter{}
		availableServiceType = map[string]struct{}{
			"delivery":     struct{}{},
			"construction": struct{}{},
			"manufacture":  struct{}{},
		}
	)

	filter.Limit = r.URL.Query().Get("limit")
	filter.Offset = r.URL.Query().Get("offset")
	filter.ServiceType = r.URL.Query()["service_type"]

	for idx, t := range filter.ServiceType {
		svType := strings.ToUpper(t)
		if _, ok := availableServiceType[svType]; !ok {
			http.Error(w, fmt.Errorf("unsuported service type %s", t).Error(), http.StatusBadRequest)

			return
		}

		filter.ServiceType[idx] = svType
	}

	convertReq, err := converter.FromStringToIntListFilter(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.List(context.Background(), convertReq.Limit, convertReq.Offset, convertReq.ServiceType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListForUser(w http.ResponseWriter, r *http.Request) {

	var (
		filter = apiModel.UserListFilter{}
	)

	filter.Limit = r.URL.Query().Get("limit")
	filter.Offset = r.URL.Query().Get("offset")
	filter.Username = r.URL.Query().Get("username")

	convertReq, err := converter.FromStringToIntUserListFilter(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.validator.Struct(convertReq); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.UserList(context.Background(), convertReq.Limit, convertReq.Offset, convertReq.Username)

	if err != nil {
		if errors.Is(err, stderrs.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(err.Error()))

		return
	}

	if len(resp) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {

	var (
		req  = apiModel.StatusRequest{}
		vars = mux.Vars(r)
	)

	tenderUUID, err := uuid.Parse(vars["tenderId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	req.TenderID = tenderUUID
	req.Username = r.URL.Query().Get("username")

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.Status(context.Background(), req.TenderID, req.Username)
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

func (h *Handler) EditStatus(w http.ResponseWriter, r *http.Request) {

	var (
		req  = apiModel.EditStatusRequest{}
		vars = mux.Vars(r)
	)

	tenderUUID, err := uuid.Parse(vars["tenderId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	req.TenderID = tenderUUID
	req.Username = r.URL.Query().Get("username")
	req.Status = strings.ToUpper(r.URL.Query().Get("status"))

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.EditStatus(context.Background(), req.TenderID, req.Username, req.Status)
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

func (h *Handler) EditTender(w http.ResponseWriter, r *http.Request) {

	var (
		req  = apiModel.EditTenderRequest{}
		vars = mux.Vars(r)
	)

	req.TenderID = vars["tenderId"]
	req.Username = r.URL.Query().Get("username")

	if err := json.NewDecoder(r.Body).Decode(&req.Params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	req.Params.ServiceType = strings.ToUpper(req.Params.ServiceType)

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := h.service.EditTender(context.Background(), &req)
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
