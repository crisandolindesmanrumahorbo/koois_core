package handler

import (
	"encoding/json"
	"net/http"
	// "strconv"

	"koois_core/internal/model"
	"koois_core/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users == nil {
		users = []model.User{}
	}

	writeJSON(w, http.StatusOK, users)
}

// func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		writeError(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}
//
// 	user, err := h.service.GetByID(r.Context(), id)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			writeError(w, http.StatusNotFound, "User not found")
// 			return
// 		}
// 		writeError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	writeJSON(w, http.StatusOK, user)
// }
//
// func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	var req model.CreateUserReq
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		writeError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}
//
// 	user, err := h.service.Create(r.Context(), req)
// 	if err != nil {
// 		writeError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	writeJSON(w, http.StatusCreated, user)
// }
//
// func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		writeError(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}
//
// 	var req model.UpdateUserReq
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		writeError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}
//
// 	user, err := h.service.Update(r.Context(), id, req)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			writeError(w, http.StatusNotFound, "User not found")
// 			return
// 		}
// 		writeError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	writeJSON(w, http.StatusOK, user)
// }
//
// func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.PathValue("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		writeError(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}
//
// 	if err := h.service.Delete(r.Context(), id); err != nil {
// 		if err.Error() == "user not found" {
// 			writeError(w, http.StatusNotFound, "User not found")
// 			return
// 		}
// 		writeError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
//
// 	writeJSON(w, http.StatusNoContent, nil)
// }

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
