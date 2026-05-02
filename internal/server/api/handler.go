package api

import (
	"gamesync/internal/server/middleware"
	"net/http"

	"xorm.io/xorm"
)

type Handler struct {
	engine *xorm.Engine
	opts HandlerOpts
}
type HandlerOpts struct {
	Logging bool
}
func NewHandler(opts HandlerOpts, engine *xorm.Engine) *Handler {
	return &Handler{
		engine: engine,
		opts: opts,
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
	}
	if h.opts.Logging {
		server.Handler = middleware.Logging(api)
	} else {
		server.Handler = api
	}

	return server.ListenAndServe()
}

func (h *Handler) getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
