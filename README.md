# Altus - Cloud Risk Visualizer

Altus est un outil de cartographie et d'audit de sécurité multi-cloud. Il permet de scanner des infrastructures, de détecter les vecteurs d'attaque critiques (exposition publique, chemins d'élévation de privilèges IAM) et de les visualiser sous forme de graphe interactif orienté (Blast Radius).

## Interface web Neo4j

<p align="center">
  <img src="https://github.com/user-attachments/assets/7ca7def2-f709-4463-a866-e117c84f1d27" alt="Interface Altus" width="800">
</p>

## 🚀 Architecture de l'application

L'outil est structuré en Monorepo :
* **`scanner/` (Go)** : Collecte rapide des ressources Cloud et exécution de la logique d'audit de sécurité.
* **`backend/` (FastAPI / Python)** : API intermédiaire qui reçoit les données du scanner, les injecte et les requête dans la base de données orientée graphe.
* **`frontend/` (React / Cytoscape.js)** : Interface utilisateur pour visualiser la topologie des risques et analyser le Blast Radius au clic.
* **Base de données** : Neo4j (Graph Database).

---

## 🛠️ Prérequis

Avant de lancer l'application, assurez-vous d'avoir installé :
* [Node.js](https://nodejs.org/) (v18+)
* [Python](https://www.python.org/) (v3.10+)
* [Go](https://go.dev/) (v1.20+)
* Une instance [Neo4j](https://neo4j.com/) active (locale ou Docker)

---

## 🏁 Ordre de lancement et commandes

Pour faire fonctionner Altus, suivez les étapes de lancement dans l'ordre exact ci-dessous.

### 1. Base de données (Neo4j)
Assurez-vous que votre instance Neo4j est démarrée et accessible sur ses ports par défaut (`bolt://localhost:7687`).

### 2. Backend (FastAPI)
Ouvrez un terminal dans le dossier `backend/` :
```bash
cd backend
# Activation de l'environnement virtuel (Windows)
.\venv\Scripts\activate
# (Ou sur macOS/Linux : source venv/bin/activate)

# Installation des dépendances
pip install -r requirements.txt

# Lancement du serveur API
uvicorn main:app --reload
```

### 3. Frontend (React)

Ouvrez un second terminal dans le dossier frontend/ :
```bash
cd frontend

# Installation des dépendances (génère le dossier node_modules local)
npm install

# Lancement de l'interface en mode développement
npm run dev
```

L'interface est accessible sur : http://localhost:5173


### 4. Scanner Go

Ouvrez un second terminal dans le dossier frontend/ :
```bash
cd scanner
# Exécution du scanneur et envoi automatique des données au backend
go run main.go
```

## 🛡️ Fonctionnalités d'Audit Actuelles

- Public Exposure Detection : Identifie les ressources exposées directement à la passerelle **INTERNET**.

- Privilege Escalation (PrivEsc) : Analyse les chemins de permissions IAM critiques permettant à un utilisateur d'obtenir des privilèges **AdministratorAccess** (relation **CAN_ESCALATE_TO**).

- Blast Radius Analysis : Cliquez sur un nœud dans l'interface graphique pour mettre en évidence toutes ses relations directes et indirectes de compromission.
