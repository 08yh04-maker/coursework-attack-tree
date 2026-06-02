package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateProbability(t *testing.T) {
	tree := AttackTree{
		Name: "Test",
		Root: AttackNode{
			ID:          "root",
			Name:        "Root",
			NodeType:    OR,
			Probability: 0,
			Cost:        0,
			Children: []AttackNode{
				{
					ID:          "leaf1",
					Name:        "Leaf1",
					NodeType:    LEAF,
					Probability: 0.3,
					Cost:        1000,
				},
				{
					ID:          "leaf2",
					Name:        "Leaf2",
					NodeType:    LEAF,
					Probability: 0.4,
					Cost:        2000,
				},
			},
		},
	}

	analyzer := NewAttackTreeAnalyzer(tree)
	prob := analyzer.CalculateProbability("root")

	// Допустимая погрешность для float
	expected := 0.58
	delta := 0.0001

	if prob < expected-delta || prob > expected+delta {
		t.Errorf("Expected probability %f, got %f", expected, prob)
	}
}

func TestCalculateCost(t *testing.T) {
	tree := AttackTree{
		Name: "Test",
		Root: AttackNode{
			ID:          "root",
			Name:        "Root",
			NodeType:    OR,
			Probability: 0,
			Cost:        0,
			Children: []AttackNode{
				{
					ID:          "leaf1",
					Name:        "Leaf1",
					NodeType:    LEAF,
					Probability: 0.3,
					Cost:        5000,
				},
				{
					ID:          "leaf2",
					Name:        "Leaf2",
					NodeType:    LEAF,
					Probability: 0.4,
					Cost:        3000,
				},
			},
		},
	}

	analyzer := NewAttackTreeAnalyzer(tree)
	cost := analyzer.CalculateCost("root")

	if cost != 3000 {
		t.Errorf("Expected cost 3000, got %f", cost)
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}