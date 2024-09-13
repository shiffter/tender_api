package bids

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"tender/internal/repository"
	serviceModel "tender/internal/service/model"
)

type repo struct {
	db *pgxpool.Pool
}

func NewBidsRepos(pool *pgxpool.Pool) repository.BidsRepos {
	return &repo{db: pool}
}

func (r *repo) Create(ctx context.Context, bid *serviceModel.CreateBidsRequest) (*serviceModel.CreateBidsResponse, error) {
	var (
		sql = `WITH org_at_usr as (
				SELECT organization_id as org_id FROM organization_responsible WHERE user_id = $1)
			INSERT INTO proposals (tender_id, author_id, organization_id, name, description)
			VALUES ($2, $3, (select org_id from org_at_usr limit 1), $4, $5)
			RETURNING id, author_id, name, description, status, version, created_at`
		resp = serviceModel.CreateBidsResponse{}
	)

	err := r.db.QueryRow(ctx, sql, bid.AuthorID, bid.TenderID, bid.AuthorID, bid.Name, bid.Description).
		Scan(&resp.ID, &resp.AuthorID, &resp.Name, &resp.Description, &resp.Status, &resp.Version, &resp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
