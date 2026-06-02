package main

type NodeType string

const (
	AND  NodeType = "AND"
	OR   NodeType = "OR"
	LEAF NodeType = "LEAF"
)

type AttackNode struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	NodeType    NodeType    `json:"node_type"`
	Probability float64     `json:"probability,omitempty"`
	Cost        float64     `json:"cost,omitempty"`
	Children    []AttackNode `json:"children,omitempty"`
}

type AttackTree struct {
	Name string     `json:"name"`
	Root AttackNode `json:"root"`
}

type RiskMetrics struct {
	TotalProbability float64   `json:"total_probability"`
	TotalCost        float64   `json:"total_cost"`
	CriticalNodes    []string  `json:"critical_nodes"`
	AttackPaths      [][]string `json:"attack_paths"`
}