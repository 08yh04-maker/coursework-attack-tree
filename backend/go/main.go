package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tree AttackTree
	if err := json.NewDecoder(r.Body).Decode(&tree); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	analyzer := NewAttackTreeAnalyzer(tree)
	metrics := analyzer.GetRiskMetrics()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	http.HandleFunc("/analyze", enableCORS(analyzeHandler))
	http.HandleFunc("/health", enableCORS(healthHandler))

	log.Println("Go Risk Calculator server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}