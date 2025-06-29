package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suraj/GoGoNotes/models"
	"github.com/suraj/GoGoNotes/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteHandler struct {
	model     *models.NoteModel
	SecretKey []byte
}

func NewNoteHandler(noteModel *models.NoteModel, secretKey []byte) *NoteHandler {
	return &NoteHandler{
		model:     noteModel,
		SecretKey: secretKey,
	}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	note, err := h.model.Create(userID, input.Title, input.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Note created successfully",
		"note":    note,
	})
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
			"notes":   []interface{}{},
		})
		return
	}

	notes, err := h.model.GetAll(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch notes: " + err.Error(),
			"notes":   []interface{}{},
		})
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Notes fetched successfully",
		"notes":   notes,
	})
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	params := mux.Vars(r)
	noteID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid note ID",
		})
		return
	}

	note, err := h.model.GetByID(noteID, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Note not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Note fetched successfully",
		"note":    note,
	})
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	params := mux.Vars(r)
	noteID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid note ID",
		})
		return
	}

	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	note, err := h.model.Update(noteID, userID, input.Title, input.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to update note: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Note updated successfully",
		"note":    note,
	})
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	params := mux.Vars(r)
	noteID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid note ID",
		})
		return
	}

	err = h.model.Delete(noteID, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to delete note: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Note deleted successfully",
	})
}
