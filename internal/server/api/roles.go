package api

import (
	"encoding/json"
	"gamesync/internal/model"
	"net/http"
)


type roleCreator struct {
	Name string `json:"name"`
}
func (h *Handler) createRole(w http.ResponseWriter, r *http.Request) {
	c := roleCreator{}
	decoder := json.NewDecoder(r.Body)
	if decoder == nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if err := decoder.Decode(&c); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	modelRole := model.Role{
		RoleName: c.Name,
	}
	if err := h.q.Role.WithContext(r.Context()).Create(&modelRole); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
