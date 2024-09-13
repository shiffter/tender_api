package model

import "time"

type CreateBidRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderID    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorID    string `json:"authorId"`
}

type CreateBidResp struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorID   string    `json:"authorId"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}
