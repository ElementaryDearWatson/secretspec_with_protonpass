package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	secretspec "github.com/cachix/secretspec/secretspec-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Resolve secrets using SecretSpec builder
	resolved, err := secretspec.New().Load()
	var apiKey string
	if err != nil {
		log.Printf("Warning: SecretSpec resolution failed, falling back to env: %v", err)
		apiKey = os.Getenv("API_KEY")
	} else {
		defer resolved.Close()
		apiKey = resolved.Secrets["API_KEY"].Get()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go app running with SecretSpec and Proton Pass!")
	})

	http.HandleFunc("/api-data", func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" {
			http.Error(w, "API Key missing", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "SecretSpec loaded API Key successfully (length: %d)", len(apiKey))
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
