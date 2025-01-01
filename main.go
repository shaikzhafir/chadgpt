package main

import (
	"chat-backend/conversor/anthropic"
	"chat-backend/conversor/openai"
	"chat-backend/internal/handlers"
	"log"
	"net/http"
	"os"
)

type message struct {
	Message string `json:"message"`
}

func main() {
	dsKey := os.Getenv("DEEPSEEK_KEY")
	oapiKey := os.Getenv("OPENAI_KEY")
	anthKey := os.Getenv("ANTHROPIC_KEY")
	server := http.NewServeMux()
	deepSeekConversor := openai.NewConversor(dsKey, "https://api.deepseek.com", "deepseek-chat")
	oapiConversor := openai.NewConversor(oapiKey, "", "gpt-4o")
	anthropicConversor := anthropic.New(anthKey)

	// css and js files
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// init handlers and routes
	oapihandler := handlers.NewChatHandler(oapiConversor)
	server.Handle("/chat-openai", oapihandler)
	deepseekHandler := handlers.NewChatHandler(deepSeekConversor)
	server.Handle("/chat-deepseek", deepseekHandler)
	anthropicHandler := handlers.NewChatHandler(anthropicConversor)
	server.Handle("/chat-anthropic", anthropicHandler)

	// catch all for html files
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", server)
}
