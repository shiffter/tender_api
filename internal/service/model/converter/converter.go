package converter

import (
	"github.com/google/uuid"
	apiModel "tender/internal/handler/model"
	serviceModel "tender/internal/service/model"
	"time"
)

func FromApiToServiceCreateReq(req *apiModel.CreateRequest) *serviceModel.CreateTenderRequest {
	return &serviceModel.CreateTenderRequest{
		Name:           req.Name,
		Description:    req.Description,
		ServiceType:    req.ServiceType,
		OrganizationID: req.OrganizationID,
		Creator:        req.Creator,
	}
}

func FromServiceToApi(tender *serviceModel.Tender) *apiModel.Tender {
	return &apiModel.Tender{
		ID:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.ServiceType,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}
}

func FromApiToServiceEditTender(req *apiModel.EditTenderRequest) (*serviceModel.EditTenderParams, error) {
	tenderID, err := uuid.Parse(req.TenderID)
	if err != nil {
		return nil, err
	}

	return &serviceModel.EditTenderParams{
		TenderID:    tenderID,
		Username:    req.Username,
		Name:        req.Params.Name,
		Description: req.Params.Description,
		ServiceType: req.Params.ServiceType,
	}, nil
}

func FromServiceToApiBidsResp(serviceResp *serviceModel.CreateBidsResponse) *apiModel.CreateBidResp {
	return &apiModel.CreateBidResp{
		ID:         serviceResp.ID.String(),
		Name:       serviceResp.Name,
		Status:     serviceResp.Status,
		AuthorType: serviceResp.AuthorType,
		AuthorID:   serviceResp.AuthorID.String(),
		Version:    0,
		CreatedAt:  time.Time{},
	}

}
