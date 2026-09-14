package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("🔴 WARNING: API_KEY is missing! Run with 'secretspec run --' to inject secrets.")
	} else {
		log.Printf("🟢 SUCCESS: API_KEY successfully resolved! (Key Length: %d)", len(apiKey))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go app running with SecretSpec and Proton Pass!")
	})

	http.HandleFunc("/api-data", func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" {
			http.Error(w, "❌ Error 500: API Key missing from environment", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "✅ SecretSpec active: API Key loaded successfully (Length: %d)", len(apiKey))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
