package handlers

import (
	"chat-backend/conversor"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ChatHandler struct {
	conversor conversor.Conversor
}

type message struct {
	Message string `json:"message"`
}

func NewChatHandler(c conversor.Conversor) *ChatHandler {
	return &ChatHandler{conversor: c}
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := getOrCreateSessionID(r)

	m, err := parseMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, sessionID, err := h.conversor.Ask(r.Context(), sessionID, m.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("session-id", sessionID)
	w.Write([]byte(resp))
}

func getOrCreateSessionID(r *http.Request) string {
	sessionID := r.Header.Get("session-id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	return sessionID
}

func parseMessage(r *http.Request) (*message, error) {
	var m message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	if m.Message == "" {
		return nil, fmt.Errorf("message is empty")
	}

	return &m, nil
}
