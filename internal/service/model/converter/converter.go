package converter

import (
	apiModel "tender/internal/handler/model"
	serviceModel "tender/internal/service/model"
)

func FromApiToServiceCreateReq(req *apiModel.CreateRequest) *serviceModel.CreateRequest {
	return &serviceModel.CreateRequest{
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
