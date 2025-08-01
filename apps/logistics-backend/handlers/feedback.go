package handlers

import (
	"encoding/json"
	"fmt"
	"logistics-backend/internal/domain/feedback"
	usecase "logistics-backend/internal/usecase/feedback"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	FH *usecase.UseCase
}

func NewFeedbackHandler(fh *usecase.UseCase) *FeedbackHandler {
	return &FeedbackHandler{FH: fh}
}

// CreateFeedback godoc
// @Summary Create a new feedback
// @Security JWT
// @Description Create a new feedback with order_id, customer_id, etc.
// @Tags feedbacks
// @Accept  json
// @Produce  json
// @Param user body feedback.CreateFeedbackRequest true "User Input"
// @Success 201 {object} feedback.Feedback
// @Failure 400 {string} handlers.ErrorResponse "Invalid request"
// @Failure 500 {string} handlers.ErrorResponse "Failed to create feedback"
// @Router /feedbacks/create [post]
func (fh *FeedbackHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var req *feedback.CreateFeedbackRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	f := req.ToFeedback()

	if err := fh.FH.CreateFeedback(r.Context(), f); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not create feedback", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":           f.ID,
		"order_id":     f.OrderID,
		"customer_id":  f.CustomerID,
		"rating":       f.Rating,
		"comments":     f.Comments,
		"submitted_at": f.SubmittedAt,
	})
}

// GetFeedbackByID godoc
// @Summary Get feedback by ID
// @Security JWT
// @Description Retrieve a feedback by their ID
// @Tags feedbacks
// @Produce  json
// @Param id path string true "Feedback ID"
// @Success 200 {object} feedback.Feedback
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Feedback not found"
// @Router /feedbacks/{id} [get]
func (fh *FeedbackHandler) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	feedbackID, err := uuid.Parse(idStr)
	fmt.Println("parsed id:", feedbackID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid feedback ID", nil)
		return
	}

	f, err := fh.FH.GetFeedbackByID(r.Context(), feedbackID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No feedback found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

// ListFeedbacks godoc
// @Summary List all feedbacks
// @Security JWT
// @Description Get a list of all feedbacks
// @Tags feedbacks
// @Produce  json
// @Success 200 {array} feedback.Feedback
// @Router /feedbacks/all_feedbacks [get]
func (fh *FeedbackHandler) ListFeedback(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := fh.FH.ListFeedback(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch feedbacks", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}
