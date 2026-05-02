package api

import (
	"net/http"

	"xorm.io/xorm"
)

type Handler struct {
	engine *xorm.Engine
}

func NewHandler(engine *xorm.Engine) *Handler {
	return &Handler{
		engine: engine,
	}
}

func (h *Handler) Serve() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", h.getHealth)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", v1))

	server := http.Server{
		Addr: ":8080",
		Handler: api,
	}
	return server.ListenAndServe()
}

func (h *Handler) getHealth(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
