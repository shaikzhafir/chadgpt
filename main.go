package main

import (
	"chat-backend/conversor"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type message struct {
	Message string `json:"message"`
}

func main() {
	authKey := os.Getenv("AUTH_KEY")
	server := http.NewServeMux()
	conversor := conversor.NewConversor(authKey)
	server.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
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
		resp, err := conversor.Ask(r.Context(), m.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(resp))
	})
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", server)

}
