package api

import (
	"gamesync/internal/model"
	"gamesync/internal/query"
	"gamesync/internal/server/middleware"
	"log/slog"
	"net/http"
)

type Handler struct {
	q *query.Query
	opts HandlerOpts
}
type HandlerOpts struct {
	Logging bool
}
func NewHandler(opts HandlerOpts, query *query.Query) *Handler {
	return &Handler{
		q: query,
		opts: opts,
	}
}

func (h *Handler) Serve() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", h.getHealth)
	router.HandleFunc("GET /users", h.getUsers)

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

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	if err := h.q.Role.WithContext(r.Context()).Create(&model.Role{RoleName: "admin"}); err != nil {
		slog.Error("creating role", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		slog.Info("created role")
		w.WriteHeader(http.StatusOK)
	}
	if err := h.q.User.WithContext(r.Context()).Create(&model.User{UserName: "test", RoleID: 1}); err != nil {
		slog.Error("creating user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		slog.Info("created user")
		w.WriteHeader(http.StatusOK)
	}
}
