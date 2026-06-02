from pydantic import BaseModel
from typing import List, Optional
from enum import Enum

class NodeType(str, Enum):
    AND = "AND"
    OR = "OR"
    LEAF = "LEAF"

class AttackNode(BaseModel):
    id: str
    name: str
    node_type: NodeType
    probability: Optional[float] = None      # для листьев
    cost: Optional[float] = None             # для листьев
    children: List['AttackNode'] = []

class AttackTree(BaseModel):
    name: str
    root: AttackNode

class RiskMetrics(BaseModel):
    total_probability: float
    total_cost: float
    critical_nodes: List[str]
    attack_paths: List[List[str]]

AttackNode.model_rebuild()