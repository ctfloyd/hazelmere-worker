package handler

import "github.com/go-chi/chi/v5"

type ApiVersion int

const (
	_ ApiVersion = iota
	ApiVersionV1
)

type WorkerHandler interface {
	RegisterRoutes(mux *chi.Mux, version ApiVersion)
}
