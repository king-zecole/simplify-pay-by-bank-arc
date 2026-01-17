// Bank Verification Microservice - Go Implementation
// Clean Architecture: Infrastructure Layer

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Domain Entities
type User struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	IBAN     string `json:"iban"`
}

type Bank struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type VerificationRequest struct {
	User User `json:"user"`
	Bank Bank `json:"bank"`
}

type VerificationResponse struct {
	Success        bool                   `json:"success"`
	Message        string                 `json:"message"`
	Timestamp      string                 `json:"timestamp"`
	VerificationID string                 `json:"verificationId,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Generate random verification ID
func generateVerificationID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 9)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("VRF-%d-%s", time.Now().Unix(), string(b))
}

// Mask IBAN for response
func maskIBAN(iban string) string {
	if len(iban) < 8 {
		return iban
	}
	return iban[:4] + strings.Repeat("*", len(iban)-8) + iban[len(iban)-4:]
}

// Validate IBAN format (basic validation)
func validateIBAN(iban string) bool {
	// Remove spaces and convert to uppercase
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))
	// Basic length check (IBANs are typically 15-34 characters)
	return len(iban) >= 15 && len(iban) <= 34
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:    "healthy",
		Service:   "bank-verification-api",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Verification handler
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(VerificationResponse{
			Success:   false,
			Message:   "Method not allowed",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerificationResponse{
			Success:   false,
			Message:   "Invalid request body",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// Validate request
	if req.User.FullName == "" || req.User.Email == "" || req.User.IBAN == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerificationResponse{
			Success:   false,
			Message:   "Missing required user fields",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	if req.Bank.ID == "" || req.Bank.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerificationResponse{
			Success:   false,
			Message:   "Missing required bank fields",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// Validate IBAN
	if !validateIBAN(req.User.IBAN) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerificationResponse{
			Success:   false,
			Message:   "Invalid IBAN format",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// Simulate verification processing delay
	time.Sleep(1500 * time.Millisecond)

	// Generate verification response
	verificationID := generateVerificationID()

	response := VerificationResponse{
		Success:        true,
		Message:        "Bank account verified successfully",
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		VerificationID: verificationID,
		Data: map[string]interface{}{
			"accountHolder": req.User.FullName,
			"bankName":      req.Bank.Name,
			"maskedIban":    maskIBAN(req.User.IBAN),
		},
	}

	log.Printf("Verification successful: ID=%s, User=%s, Bank=%s",
		verificationID, req.User.Email, req.Bank.Name)

	json.NewEncoder(w).Encode(response)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")

	// Routes
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/api/verify", corsMiddleware(verifyHandler))

	log.Printf("Bank Verification Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
