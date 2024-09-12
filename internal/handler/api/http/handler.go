package handler

import "net/http"

type TenderHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	ListForUser(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
	EditStatus(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
}
