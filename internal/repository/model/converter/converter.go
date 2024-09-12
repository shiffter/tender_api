package converter

import (
	repoModel "tender/internal/repository/model"
	serviceModel "tender/internal/service/model"
)

func ToServiceTenderFromRepo(tender *repoModel.Tender) *serviceModel.Tender {
	return &serviceModel.Tender{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		Status:         tender.Status,
		Creator:        tender.Creator,
		OrganizationId: tender.OrganizationID,
		ServiceType:    tender.ServiceType,
		Version:        tender.Version,
		CreatedAt:      tender.CreatedAt,
	}
}

func ToServiceUserFromRepo(user *repoModel.User) *serviceModel.User {
	return &serviceModel.User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
