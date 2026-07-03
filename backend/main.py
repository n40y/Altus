# backend/main.py

from pydantic import BaseModel
from typing import List
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from models import Node, Relation
from graph import create_node, create_relation
from collectors.aws_mock import collect
from rules import check_security_rules
from queries import get_graph


app = FastAPI(title="Altus CloudGraph API")

# Autorise le frontend React à communiquer avec l'API
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
def root():
    return {"status": "online", "app": "Altus Backend"}


@app.get("/scan")
def trigger_scan():
    # Exécute l'analyse des règles sur l'état actuel de la base Neo4j
    findings = check_security_rules()
    return {"status": "finished", "findings": findings}


@app.get("/graph")
def graph():
    return {"nodes": get_graph()}


# Schémas de validation pour les données envoyées par le scanner Go
class GoNode(BaseModel):
    name: str
    type: str

class GoRelation(BaseModel):
    source: str
    target: str
    relation: str

class ScanPayload(BaseModel):
    nodes: List[GoNode]
    relations: List[GoRelation]


@app.post("/scanner/submit")
def receive_scan(payload: ScanPayload):
    # 1. Injection des nœuds transmis par le scanner Go
    for n in payload.nodes:
        create_node(Node(name=n.name, type=n.type))
        
    # 2. Injection des relations correspondantes
    for r in payload.relations:
        create_relation(Relation(source=r.source, target=r.target, relation=r.relation))
        
    # 3. Exécution immédiate du moteur de règles Python après l'import
    active_findings = check_security_rules()
        
    return {
        "status": "success", 
        "inserted_nodes": len(payload.nodes), 
        "inserted_relations": len(payload.relations),
        "findings_detected": len(active_findings)
    }