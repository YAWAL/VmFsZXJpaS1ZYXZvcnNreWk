package router

import (
	"net/http"

	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/database"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/handlers"
	"github.com/YAWAL/VmFsZXJpaS1ZYXZvcnNreWk/src/worker"

	"github.com/gorilla/mux"
)

func New(er database.Repository, w worker.Worker) (r *mux.Router) {
	r = mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/fetcher", handlers.SaveURLData(er, w)).Methods(http.MethodPost)
	api.HandleFunc("/fetcher", handlers.GetAllURLData(er)).Methods(http.MethodGet)
	api.HandleFunc("/fetcher/{id}", handlers.DeleteURLData(er)).Methods(http.MethodDelete)
	api.HandleFunc("/fetcher/{id}/history", handlers.GetDownloadHistoriesByURLDataID(er)).Methods(http.MethodGet)
	return r
}
