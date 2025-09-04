package handlers

import (
	"encoding/json"
	"logistics-backend/internal/domain/invite"
	usecase "logistics-backend/internal/usecase/invite"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InviteHandler struct {
	UC *usecase.UseCase
}

func NewInviteHandler(uc *usecase.UseCase) *InviteHandler {
	return &InviteHandler{UC: uc}
}

// CreateMember godoc
// @Summary Create a new invite
// @Description Create a new invite for a user
// @Tags Invites
// @Accept json
// @Produce json
// @Param invite body invite.CreateInviteRequest true "Invite payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /invites/create [post]
func (h *InviteHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	var req invite.CreateInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	i := req.ToInvite()

	if err := h.UC.InviteMember(r.Context(), i); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create invite", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":         i.ID,
		"email":      i.Email,
		"role":       i.Role,
		"token":      i.Token,
		"expires_at": i.Token,
		"invited_by": i.InvitedBy,
	})
}

// GetMemberByToken godoc
// @Summary Get invite by token
// @Description Fetch an invite using a token
// @Tags Invites
// @Produce json
// @Param token query string true "Invite token"
// @Success 200 {object} invite.Invite
// @Failure 400 {string} handlers.ErrorResponse "Invalid Token"
// @Failure 404 {string} handlers.ErrorResponse "Invite not found"
// @Router /invites/by-token [get]
func (h *InviteHandler) GetMemberByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSONError(w, http.StatusBadRequest, "missing token", nil)
		return
	}

	invite, err := h.UC.GetMemberByToken(r.Context(), token)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Invite not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invite)
}

// ListPendingMembers godoc
// @Summary List all pending invites
// @Description Get all invites that are pending
// @Tags Invites
// @Produce json
// @Success 200 {array} invite.Invite
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /invites/all_invites [get]
func (h *InviteHandler) ListPendingMembers(w http.ResponseWriter, r *http.Request) {
	invites, err := h.UC.ListPendingMembers(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch invites", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invites)
}

// DeleteMember godoc
// @Summary Delete an invite
// @Description Delete an invite by ID
// @Tags Invites
// @Produce json
// @Param id path string true "Invite ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid invite ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /invites/{id} [delete]
func (h *InviteHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	inviteID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	if err := h.UC.DeleteMember(r.Context(), inviteID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete invite", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "New Invite deleted"})
}
