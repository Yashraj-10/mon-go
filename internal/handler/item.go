package handler

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"mon-go/internal/model"
	"mon-go/internal/store"
)

// ItemHandler handles HTTP for items. Depends on ItemStore for easy testing/swapping.
// Store may be nil when the database is unavailable; handlers return 503 in that case.
type ItemHandler struct {
	Store *store.ItemStore
}

func (h *ItemHandler) requireStore(w http.ResponseWriter) bool {
	if h.Store == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return false
	}
	return true
}

// CreateItem creates a new item (POST/GET /items.create).
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	var input model.ItemCreate
	if r.Method == http.MethodGet {
		input.Name = r.URL.Query().Get("name")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
	}
	if input.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	item, err := h.Store.Create(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

// GetItem returns one item by ID (POST/GET /items.get?id=...).
func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	item, err := h.Store.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// ListItems returns a list of items (POST/GET /items.list).
func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	items, err := h.Store.List(r.Context(), 0)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

// DeleteItem deletes an item by ID (POST/GET /items.delete?id=...).
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	ok, err := h.Store.DeleteByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
