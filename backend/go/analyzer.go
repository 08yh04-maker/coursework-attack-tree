package main

import (
	"math"
	"sync"
)

type AttackTreeAnalyzer struct {
	tree  AttackTree
	cache map[string]float64
	mu    sync.RWMutex
}

func NewAttackTreeAnalyzer(tree AttackTree) *AttackTreeAnalyzer {
	return &AttackTreeAnalyzer{
		tree:  tree,
		cache: make(map[string]float64),
	}
}

func (a *AttackTreeAnalyzer) findNode(nodeID string, node AttackNode) *AttackNode {
	if node.ID == nodeID {
		return &node
	}
	for _, child := range node.Children {
		if found := a.findNode(nodeID, child); found != nil {
			return found
		}
	}
	return nil
}

func (a *AttackTreeAnalyzer) CalculateProbability(nodeID string) float64 {
	a.mu.RLock()
	if val, exists := a.cache[nodeID]; exists {
		a.mu.RUnlock()
		return val
	}
	a.mu.RUnlock()

	node := a.findNode(nodeID, a.tree.Root)
	if node == nil {
		return 0.0
	}

	var result float64

	if node.NodeType == LEAF {
		result = node.Probability
	} else if node.NodeType == AND {
		result = 1.0
		for _, child := range node.Children {
			result *= a.CalculateProbability(child.ID)
		}
	} else if node.NodeType == OR {
		result = 1.0
		for _, child := range node.Children {
			result *= (1 - a.CalculateProbability(child.ID))
		}
		result = 1 - result
	}

	a.mu.Lock()
	a.cache[nodeID] = result
	a.mu.Unlock()

	return result
}

func (a *AttackTreeAnalyzer) CalculateCost(nodeID string) float64 {
	node := a.findNode(nodeID, a.tree.Root)
	if node == nil {
		return 0.0
	}

	if node.NodeType == LEAF {
		return node.Cost
	}

	if node.NodeType == AND {
		total := 0.0
		for _, child := range node.Children {
			total += a.CalculateCost(child.ID)
		}
		return total
	}

	if node.NodeType == OR {
		minCost := math.MaxFloat64
		for _, child := range node.Children {
			cost := a.CalculateCost(child.ID)
			if cost < minCost {
				minCost = cost
			}
		}
		if minCost == math.MaxFloat64 {
			return 0.0
		}
		return minCost
	}

	return 0.0
}

func (a *AttackTreeAnalyzer) GetRiskMetrics() RiskMetrics {
	rootID := a.tree.Root.ID

	totalProb := a.CalculateProbability(rootID)
	totalCost := a.CalculateCost(rootID)

	// Критические узлы (условно: вероятность > 0.7)
	criticalNodes := make([]string, 0)
	var collectCritical func(node AttackNode)
	collectCritical = func(node AttackNode) {
		prob := a.CalculateProbability(node.ID)
		if prob > 0.7 {
			criticalNodes = append(criticalNodes, node.ID)
		}
		for _, child := range node.Children {
			collectCritical(child)
		}
	}
	collectCritical(a.tree.Root)

	return RiskMetrics{
		TotalProbability: math.Round(totalProb*10000) / 10000,
		TotalCost:        math.Round(totalCost*100) / 100,
		CriticalNodes:    criticalNodes,
		AttackPaths:      a.findAllPaths(),
	}
}

func (a *AttackTreeAnalyzer) findAllPaths() [][]string {
	paths := make([][]string, 0)
	var dfs func(node AttackNode, currentPath []string)
	dfs = func(node AttackNode, currentPath []string) {
		newPath := append(currentPath, node.ID)
		if len(node.Children) == 0 {
			pathCopy := make([]string, len(newPath))
			copy(pathCopy, newPath)
			paths = append(paths, pathCopy)
			return
		}
		for _, child := range node.Children {
			dfs(child, newPath)
		}
	}
	dfs(a.tree.Root, []string{})
	return paths
}