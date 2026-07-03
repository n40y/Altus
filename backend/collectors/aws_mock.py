# backend/collectors/aws_mock.py

from models import Node, Relation
from graph import create_node, create_relation


def collect():
    # Simulation de ressources cloud
    u = Node(name="user-test-admin", type="IAM_USER")
    p = Node(name="AdminPolicy", type="IAM_POLICY")
    b = Node(name="altus-prod-database-backup", type="S3_BUCKET")
    
    create_node(u)
    create_node(p)
    create_node(b)
    
    create_relation(Relation(source="user-test-admin", target="AdminPolicy", relation="HAS_POLICY"))
    create_relation(Relation(source="AdminPolicy", target="altus-prod-database-backup", relation="ACCESS_DATA"))