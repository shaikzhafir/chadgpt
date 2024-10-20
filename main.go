package main

import (
	"chat-backend/conversor/anthropic"
	"chat-backend/conversor/openai"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type message struct {
	Message string `json:"message"`
}

func main() {
	authKey := os.Getenv("AUTH_KEY")
	anthKey := os.Getenv("ANTHROPIC_KEY")
	server := http.NewServeMux()
	conversor := openai.NewConversor(authKey)
	anthropicConversor := anthropic.New(anthKey)

	// css and js files
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// routes
	server.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("session-id")
		if sessionID == "" {
			sessionID = uuid.New().String()
		}
		// extract user message from json body
		var m message
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if m.Message == "" {
			http.Error(w, "message is empty", http.StatusBadRequest)
			return
		}
		resp, sessionID, err := conversor.Ask(r.Context(), sessionID, m.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("session-id", sessionID)
		w.Write([]byte(resp))
	})

	server.HandleFunc("/chat-anthropic", func(w http.ResponseWriter, r *http.Request) {
		// check if req has session-id. if have, then we use it to reference context
		sessionID := r.Header.Get("session-id")
		log.Printf("sessionID: %s", sessionID)
		if sessionID == "" {
			sessionID = uuid.New().String()
		}
		// extract user message from json body
		var m message
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing form from user input %s", err.Error()), http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if m.Message == "" {
			http.Error(w, "message is empty", http.StatusBadRequest)
			return
		}
		resp, sessionID, err := anthropicConversor.Ask(r.Context(), sessionID, m.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("session-id", sessionID)
		w.Write([]byte(resp))
	})

	// catch all for html files
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", server)

}
