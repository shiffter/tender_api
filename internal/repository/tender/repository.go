package tender

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"strings"
	repoModel "tender/internal/repository/model"
	"tender/internal/repository/model/converter"
	serviceModel "tender/internal/service/model"

	"github.com/jackc/pgx/v4/pgxpool"

	"tender/internal/repository"
)

type repo struct {
	db *pgxpool.Pool
}

func NewTenderRepos(pool *pgxpool.Pool) repository.TenderRepos {
	return &repo{db: pool}
}

func (r *repo) Get(ctx context.Context, tenderID uuid.UUID) (*serviceModel.Tender, error) {
	var (
		sql = `SELECT id, name, description, status,
				creator_username, organization_id, service_type, version, created_at 
				FROM tender 
				WHERE id=$1`
		tender = repoModel.Tender{}
	)

	err := r.db.QueryRow(ctx, sql, tenderID).
		Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.Creator,
			&tender.OrganizationID, &tender.ServiceType, &tender.Version, &tender.CreatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceTenderFromRepo(&tender), nil
}

func (r *repo) Create(ctx context.Context, params *serviceModel.CreateRequest) (*serviceModel.Tender, error) {
	var (
		sql = `INSERT INTO 
			tender(name, description, status, organization_id, creator_username, service_type)
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id, name, description, status, service_type, version, created_at`
		tender = repoModel.Tender{}
	)

	err := r.db.QueryRow(ctx, sql, params.Name, params.Description, "CREATED",
		params.OrganizationID, params.Creator, params.ServiceType).
		Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status,
			&tender.ServiceType, &tender.Version, &tender.CreatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceTenderFromRepo(&tender), nil
}

func (r *repo) OrganizationRightsForUser(ctx context.Context, userName string, orgUUID uuid.UUID) (*uuid.UUID, error) {
	var (
		sql = `SELECT organization_responsible.id FROM organization_responsible 
		JOIN public.employee e
		ON organization_responsible.user_id = e.id
		WHERE e.username = $1
		AND organization_responsible.organization_id = $2`
		respID = uuid.UUID{}
	)

	err := r.db.QueryRow(ctx, sql, userName, orgUUID).Scan(&respID)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, fmt.Errorf("user %s dont have rights from the organization", userName)
	case err != nil:
		return nil, err
	}

	return &respID, nil
}

func (r *repo) List(ctx context.Context, limit, offset int32, serviceTypes []string) ([]*serviceModel.Tender, error) {
	var (
		sql, args = buildListQuery(limit, offset, serviceTypes)
		tenders   = make([]*serviceModel.Tender, 0)
	)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tmp serviceModel.Tender
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Status, &tmp.ServiceType, &tmp.Version, &tmp.CreatedAt)
		if err != nil {
			return nil, err
		}

		tenders = append(tenders, &tmp)
	}

	return tenders, nil
}

func (r *repo) UserList(ctx context.Context, limit, offset int32, username string) ([]*serviceModel.Tender, error) {
	var (
		sql = `SELECT id, name, description, status,
				service_type, version, created_at
				FROM tender WHERE creator_username = $1
				ORDER BY name
				LIMIT $2
				OFFSET $3`
		tenders = make([]*serviceModel.Tender, 0)
	)

	rows, err := r.db.Query(ctx, sql, username, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tmp serviceModel.Tender
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Status, &tmp.ServiceType, &tmp.Version, &tmp.CreatedAt)
		if err != nil {
			return nil, err
		}

		tenders = append(tenders, &tmp)
	}

	return tenders, nil
}

func (r *repo) EditStatus(ctx context.Context, status string, tenderID uuid.UUID) (*serviceModel.Tender, error) {
	var (
		sql = `UPDATE tender SET status = $1 WHERE id = $2 
		RETURNING id, name, description, status, creator_username, organization_id, service_type, version, created_at`
		updTender repoModel.Tender
	)

	err := r.db.QueryRow(ctx, sql, status, tenderID).
		Scan(&updTender.ID, &updTender.Name, &updTender.Description,
			&updTender.Status, &updTender.Creator, &updTender.OrganizationID,
			&updTender.ServiceType, &updTender.Version, &updTender.CreatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToServiceTenderFromRepo(&updTender), nil
}

func buildListQuery(limit, offset int32, serviceType []string) (string, []interface{}) {
	var (
		queryBuilder strings.Builder
		args         []interface{}
	)

	if len(serviceType) == 0 {
		queryBuilder.WriteString(
			`SELECT id, name, description, status,
				service_type, version, created_at
				FROM tender 
				ORDER BY name
				LIMIT $1
				OFFSET $2`)
		args = append(args, limit, offset)
	} else {
		queryBuilder.WriteString(
			`SELECT id, name, description, status,
				service_type, version, created_at
				FROM tender WHERE service_type = ANY($1)
				ORDER BY name
				LIMIT $2
				OFFSET $3`)
		args = append(args, serviceType, limit, offset)
	}

	return queryBuilder.String(), args
}
