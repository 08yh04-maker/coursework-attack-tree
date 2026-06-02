from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from models import AttackTree, RiskMetrics
from analyzer import AttackTreeAnalyzer

app = FastAPI(title="Attack Tree Analyzer", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def root():
    return {"message": "Attack Tree Analyzer API", "status": "ok"}

@app.post("/analyze", response_model=RiskMetrics)
def analyze_tree(tree: AttackTree):
    try:
        analyzer = AttackTreeAnalyzer(tree)
        metrics = analyzer.get_risk_metrics()
        return metrics
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))

@app.get("/health")
def health():
    return {"status": "ok"}