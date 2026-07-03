# backend/database.py

from neo4j import GraphDatabase


class Database:
    
    def __init__(self) -> None:
        self.driver = GraphDatabase.driver(
            "bolt://localhost:7687",
            auth=("neo4j", "password123")
        )
        
    def close(self):
        self.driver.close()
        
    def query(self, query, params={}):
        with self.driver.session() as session:
            result = session.run(query, params)
            return list(result)