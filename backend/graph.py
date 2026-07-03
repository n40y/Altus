# backend/graph.py

from database import Database


db = Database()

def create_node(node):
    query = """
    MERGE (n:Resource {
        name: $name,
        type: $type
    })
    """
    db.query(query, {"name": node.name, "type": node.type})
   
def create_relation(rel):
    query = """
    MATCH (a:Resource {name: $source}), (b:Resource {name: $target})
    MERGE (a)-[:RELATION {type: $relation}]->(b)
    """
    db.query(query, {
        "source": rel.source,
        "target": rel.target, 
        "relation": rel.relation
    })