# backend/models.py

from dataclasses import dataclass

@dataclass
class Node:
    name: str
    type: str
    

@dataclass
class Relation:
    source: str
    target: str
    relation: str