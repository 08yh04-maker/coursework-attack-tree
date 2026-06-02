import pytest
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"

def test_analyze():
    test_tree = {
        "name": "Test Attack",
        "root": {
            "id": "root",
            "name": "Root",
            "node_type": "LEAF",
            "probability": 0.5,
            "cost": 1000,
            "children": []
        }
    }
    response = client.post("/analyze", json=test_tree)
    assert response.status_code == 200
    data = response.json()
    assert "total_probability" in data
    assert "total_cost" in data