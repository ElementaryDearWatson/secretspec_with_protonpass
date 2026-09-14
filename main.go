package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/secretspec/secretspec-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Fetch API key securely using SecretSpec with Proton Pass provider
	apiKey, err := secretspec.GetSecret("API_KEY")
	if err != nil {
		log.Printf("Warning: API_KEY not resolved via SecretSpec (Proton Pass), falling back to env: %v", err)
		apiKey = os.Getenv("API_KEY")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go app running with SecretSpec and Proton Pass!")
	})

	http.HandleFunc("/api-data", func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" {
			http.Error(w, "API Key missing", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "SecretSpec loaded API Key from Proton Pass successfully (length: %d)", len(apiKey))
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
