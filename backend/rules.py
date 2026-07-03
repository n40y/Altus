# backend/rules.py

from database import Database


db = Database()

def check_security_rules():
    findings = []
    
    # Règle 1 : Détection d'un chemin d'accès critique vers un Bucket S3
    # Cherche un nœud qui peut atteindre un S3_BUCKET via une relation intermédiaire
    query_s3 = """
    MATCH (u {type: 'IAM_USER'})-[r1]->(p {type: 'IAM_POLICY'})-[r2]->(b {type: 'S3_BUCKET'})
    RETURN u.name AS user, p.name AS policy, b.name AS bucket, r2.type AS access_type
    """
    
    results = db.query(query_s3)
    
    for record in results:
        findings.append({
            "title": "Chiffrement ou Fuite de données potentielle sur un Bucket S3",
            "severity": "CRITICAL",
            "resource": record["bucket"],
            "description": f"L'utilisateur '{record['user']}' peut accéder au bucket '{record['bucket']}' en utilisant la politique '{record['policy']}' avec les droits de type {record['access_type']}."
        })

    # Règle 2 : Détection d'une instance EC2 qui utilise un rôle trop permissif
    query_ec2 = """
    MATCH (e {type: 'EC2_INSTANCE'})-[r1]->(p {type: 'IAM_POLICY'})-[r2]->(b {type: 'S3_BUCKET'})
    RETURN e.name AS instance, b.name AS bucket
    """
    
    results_ec2 = db.query(query_ec2)
    for record in results_ec2:
        findings.append({
            "title": "Instance EC2 surexposée",
            "severity": "HIGH",
            "resource": record["instance"],
            "description": f"L'instance EC2 '{record['instance']}' possède un chemin d'accès direct vers le bucket S3 '{record['bucket']}'."
        })
        
    return findings