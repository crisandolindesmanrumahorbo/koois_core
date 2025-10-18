package handler

import (
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
