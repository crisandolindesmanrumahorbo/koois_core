package handler

import (
	"encoding/json"
	"koois_core/internal/middleware"
	"koois_core/internal/model"
	"koois_core/internal/service"
	"net/http"
	"strconv"
)

type QuizHandler struct {
	service *service.QuizService
}

func NewQuizHandler(service *service.QuizService) *QuizHandler {
	return &QuizHandler{service: service}
}

func (h *QuizHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid quiz ID")
		return
	}

	quizzes, err := h.service.GetQuizQuestions(r.Context(), id)
	if err != nil {
		if err.Error() == "Quiz not found" {
			writeError(w, http.StatusNotFound, "Quiz not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, quizzes)
}

func (h *QuizHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetClaimsFromContext(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, `{"error":"unauthorized"}`)
		return
	}
	userId := claims.Sub
	id, err := strconv.Atoi(userId)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	var req model.CreateQuizReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	quizzes, err := h.service.Create(r.Context(), req, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, quizzes)
}
