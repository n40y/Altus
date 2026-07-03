# backend/queries.py

from database import Database


db = Database()

def get_graph():
    query = """
    MATCH (n)-[r]->(m)
    RETURN
        n.name AS source,
        n.type AS source_type,
        type(r) AS relation,
        m.name AS target,
        m.type AS target_type
    """
    results = db.query(query)
    
    graph = []
    for r in results:
        graph.append({
            "source": r["source"],
            "source_type": r["source_type"],
            "relation": r["relation"],
            "target": r["target"],
            "target_type": r["target_type"]
        })
    return graph