package service

import (
)

// type userCreator struct {
// 	Name string `json:"name"`
// 	RoleID int32 `json:"role_id"`
// }
// func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
// 	c := userCreator{}
// 	decoder := json.NewDecoder(r.Body)
// 	if decoder == nil {
// 		w.WriteHeader(http.StatusNotAcceptable)
// 		return
// 	}
// 	if err := decoder.Decode(&c); err != nil {
// 		w.WriteHeader(http.StatusNotAcceptable)
// 		return
// 	}
//
// 	modelUser := model.User{
// 		UserName: c.Name,
// 		RoleID: c.RoleID,
// 	}
// 	if err := s.q.User.WithContext(r.Context()).Create(&modelUser); err != nil {
// 		slog.Error("creating user", "error", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }
