import networkx as nx
from typing import List, Tuple
from models import AttackTree, AttackNode, RiskMetrics

class AttackTreeAnalyzer:
    def __init__(self, tree: AttackTree):
        self.tree = tree
        self.graph = nx.DiGraph()
        self._build_graph()

    def _build_graph(self):
        """Строит граф из дерева атак"""
        def add_children(parent_id: str, children: List[AttackNode]):
            for child in children:
                self.graph.add_edge(parent_id, child.id)
                if child.children:
                    add_children(child.id, child.children)

        root = self.tree.root
        self.graph.add_node(root.id, 
                            name=root.name, 
                            node_type=root.node_type,
                            probability=root.probability,
                            cost=root.cost)
        if root.children:
            add_children(root.id, root.children)

    def calculate_probability(self, node_id: str) -> float:
        """Рекурсивный расчёт вероятности успеха"""
        node = self._find_node(node_id)
        if not node:
            return 0.0

        # Лист: возвращаем свою вероятность
        if node.node_type == "LEAF":
            return node.probability or 0.0

        # Собираем вероятности детей
        child_probs = []
        for child in node.children:
            child_probs.append(self.calculate_probability(child.id))

        if not child_probs:
            return 0.0

        # AND: произведение вероятностей
        if node.node_type == "AND":
            result = 1.0
            for p in child_probs:
                result *= p
            return result

        # OR: 1 - (1-p1)*(1-p2)*...
        if node.node_type == "OR":
            result = 1.0
            for p in child_probs:
                result *= (1 - p)
            return 1 - result

        return 0.0

    def calculate_cost(self, node_id: str) -> float:
        """Рекурсивный расчёт стоимости атаки"""
        node = self._find_node(node_id)
        if not node:
            return 0.0

        # Лист: возвращаем свою стоимость
        if node.node_type == "LEAF":
            return node.cost or 0.0

        # Суммируем стоимость детей (минимальный путь для OR)
        if node.node_type == "AND":
            total = 0.0
            for child in node.children:
                total += self.calculate_cost(child.id)
            return total

        if node.node_type == "OR":
            min_cost = float('inf')
            for child in node.children:
                cost = self.calculate_cost(child.id)
                if cost < min_cost:
                    min_cost = cost
            return min_cost if min_cost != float('inf') else 0.0

        return 0.0

    def find_all_paths(self) -> List[List[str]]:
        """Находит все пути от корня до листьев"""
        paths = []
        leaves = [n for n in self.graph.nodes() 
                  if self.graph.out_degree(n) == 0]
        
        for leaf in leaves:
            for path in nx.all_simple_paths(self.graph, 
                                           self.tree.root.id, 
                                           leaf):
                paths.append(path)
        return paths

    def get_risk_metrics(self) -> RiskMetrics:
        """Получает полные метрики риска"""
        root_id = self.tree.root.id
        
        total_prob = self.calculate_probability(root_id)
        total_cost = self.calculate_cost(root_id)
        
        # Критические узлы (с высокой вероятностью)
        critical_nodes = []
        for node in self.graph.nodes():
            prob = self.calculate_probability(node)
            if prob > 0.7:
                critical_nodes.append(node)
        
        attack_paths = self.find_all_paths()
        
        return RiskMetrics(
            total_probability=round(total_prob, 4),
            total_cost=round(total_cost, 2),
            critical_nodes=critical_nodes,
            attack_paths=attack_paths
        )

    def _find_node(self, node_id: str) -> AttackNode:
        """Находит узел по ID"""
        def search(node: AttackNode) -> AttackNode:
            if node.id == node_id:
                return node
            for child in node.children:
                result = search(child)
                if result:
                    return result
            return None
        return search(self.tree.root)