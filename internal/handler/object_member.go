package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"

	"mon-go/internal/model"
	"mon-go/internal/store"
)

// object_id: group/[0-9]+. member_id: group/[0-9]+ or user/[0-9]+.
var (
	objectIDRegex  = regexp.MustCompile(`^group/[0-9]+$`)
	memberIDRegex  = regexp.MustCompile(`^(group|user)/[0-9]+$`)
)

// ObjectMemberHandler handles HTTP for object-members.
type ObjectMemberHandler struct {
	Store *store.ObjectMemberStore
}

func (h *ObjectMemberHandler) requireStore(w http.ResponseWriter) bool {
	if h.Store == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return false
	}
	return true
}

func parseObjectMemberInput(r *http.Request) (objectID, memberID string, ok bool) {
	if r.Method == http.MethodGet {
		return r.URL.Query().Get("object_id"), r.URL.Query().Get("member_id"), true
	}
	var input model.ObjectMemberCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return "", "", false
	}
	return input.ObjectID, input.MemberID, true
}

// deleteInput holds either id or (object_id, member_id) for object-members.delete.
type deleteInput struct {
	ID       string `json:"id"`
	ObjectID string `json:"object_id"`
	MemberID string `json:"member_id"`
}

func parseDeleteInput(r *http.Request) deleteInput {
	if r.Method == http.MethodGet {
		return deleteInput{
			ID:       r.URL.Query().Get("id"),
			ObjectID: r.URL.Query().Get("object_id"),
			MemberID: r.URL.Query().Get("member_id"),
		}
	}
	var d deleteInput
	_ = json.NewDecoder(r.Body).Decode(&d)
	return d
}

func isDuplicateKey(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// Create adds an object-member link (POST/GET /object-members.create).
func (h *ObjectMemberHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	objectID, memberID, ok := parseObjectMemberInput(r)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "object_id and member_id required"})
		return
	}
	if !objectIDRegex.MatchString(objectID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "object_id must match group/[0-9]+"})
		return
	}
	if !memberIDRegex.MatchString(memberID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "member_id must match group/[0-9]+ or user/[0-9]+"})
		return
	}

	doc, err := h.Store.Create(r.Context(), objectID, memberID)
	if err != nil {
		if isDuplicateKey(err) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "object_id and member_id already linked"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, doc)
}

// Delete removes an object-member link (POST/GET /object-members.delete).
// Caller may provide either id (object_id:member_id) or both object_id and member_id.
func (h *ObjectMemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.requireStore(w) {
		return
	}
	d := parseDeleteInput(r)

	if d.ID != "" {
		deleted, err := h.Store.DeleteByID(r.Context(), d.ID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if !deleted {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
		return
	}

	if d.ObjectID == "" || d.MemberID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "provide id or both object_id and member_id"})
		return
	}
	if !objectIDRegex.MatchString(d.ObjectID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "object_id must match group/[0-9]+"})
		return
	}
	if !memberIDRegex.MatchString(d.MemberID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "member_id must match group/[0-9]+ or user/[0-9]+"})
		return
	}

	deleted, err := h.Store.DeleteByObjectAndMember(r.Context(), d.ObjectID, d.MemberID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
