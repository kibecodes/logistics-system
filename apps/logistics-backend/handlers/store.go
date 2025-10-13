package handlers

import (
	"encoding/json"
	"fmt"
	"logistics-backend/internal/domain/store"
	usecase "logistics-backend/internal/usecase/store"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type StoreHandler struct {
	UC *usecase.UseCase
}

func NewStoreHandler(uc *usecase.UseCase) *StoreHandler {
	return &StoreHandler{UC: uc}
}

// CreateStore godoc
// @Summary Create a new store
// @Security JWT
// @Description Creates a new store for an owner and returns the created store
// @Tags stores
// @Accept json
// @Produce json
// @Param store body store.CreateStoreRequest true "Store input"
// @Success 201 {object} store.Store
// @Failure 400 {object} handlers.ErrorResponse "Invalid request"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/create [post]
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req store.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	s := req.ToStore()

	if err := h.UC.CreateStore(r.Context(), s); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not create order", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":          s.ID,
		"owner_id":    s.OwnerID,
		"name":        s.Name,
		"slug":        s.Slug,
		"description": s.Description,
		"logo_url":    s.LogoURL,
		"banner_url":  s.BannerURL,
		"is_public":   s.IsPublic,
		"created_at":  s.CreatedAt,
		"updated_at":  s.UpdatedAt,
	})
}

// GetStore godoc
// @Summary Get store by id
// @Description Fetches a single store using UUID
// @Tags stores
// @Produce json
// @Param slug query string true "Store ID"
// @Success 200 {object} store.Store
// @Failure 400 {object} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {object} handlers.ErrorResponse "Not found"
// @Router /stores/by-id/{id} [get]
func (h *StoreHandler) GetStoreByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	s, err := h.UC.GetStoreByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Store not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// GetStore godoc
// @Summary Get store by slug
// @Description Fetches a single store by its slug
// @Tags stores
// @Produce json
// @Param slug query string true "Store slug"
// @Success 200 {object} store.Store
// @Failure 400 {object} handlers.ErrorResponse "Missing slug"
// @Failure 404 {object} handlers.ErrorResponse "Store not found"
// @Router /stores/by-slug [get]
func (h *StoreHandler) GetStoreBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		writeJSONError(w, http.StatusBadRequest, "Slug query parameter is required", nil)
		return
	}

	s, err := h.UC.GetStoreBySlug(r.Context(), slug)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Store not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// UpdateStore godoc
// @Summary Update store
// @Security JWT
// @Description Update a specific field of an existing store
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param update body store.UpdateStoreRequest true "Field and value to update"
// @Success 200 {object} map[string]string "Store updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid store ID or request body"
// @Failure 404 {object} handlers.ErrorResponse "Store not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/{id}/update [put]
func (h *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	storeID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	var req store.UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	column := strings.TrimSpace(strings.ToLower(req.Column))
	if column == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing or invalid column name", err)
		return
	}

	if err := h.UC.UpdateStore(r.Context(), storeID, column, req.Value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update store", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("store %s updated successfully", column),
	})
}

// GetPublicStores godoc
// @Summary List all public stores
// @Description Returns all stores marked as public
// @Tags stores
// @Produce json
// @Success 200 {array} store.Store
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/public [get]
func (h *StoreHandler) GetPublicStores(w http.ResponseWriter, r *http.Request) {
	stores, err := h.UC.GetPublicStores(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch stores", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

// DeleteStore godoc
// @Summary Delete a store
// @Security JWT
// @Description Deletes a store by its ID
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} map[string]string "Store deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid store ID"
// @Failure 500 {object} handlers.ErrorResponse "Failed to delete store"
// @Router /stores/{id} [delete]
func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	storeID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	if err := h.UC.DeleteStore(r.Context(), storeID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete store", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("store %s deleted", storeID),
	})
}
