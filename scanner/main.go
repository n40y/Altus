// scanner/main.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"altus/scanner/collection"
	"altus/scanner/models"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	fmt.Println("[*] Altus Scanner - Démarrage de la collecte Go réelle...")

	// 1. Charger la configuration AWS depuis l'environnement du terminal
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("[-] Erreur critique lors du chargement de la config AWS : %v", err)
	}

	// Initialisation de la structure de résultats
	data := &models.ScanResult{
		Nodes:     []models.Node{},
		Relations: []models.Relation{},
	}

	// 2. Collecte des ressources réelles sur le compte AWS cible
	fmt.Println("[*] Interrogation des API AWS (IAM, S3, EC2)...")
	collection.CollectAWSResources(cfg, data)

	// 3. Lancement des modules d'audit de sécurité sur les données réelles
	fmt.Println("[*] Lancement de l'audit d'exposition publique...")
	collection.RunSecurityAudit(data)

	fmt.Println("[*] Lancement de l'audit d'élévation de privilèges (PrivEsc)...")
	collection.RunPrivEscAudit(data)

	fmt.Printf("[+] %d nœuds et %d relations identifiés.\n", len(data.Nodes), len(data.Relations))

	// 4. Sérialisation des données d'infrastructure et d'audit
	fmt.Println("[*] Envoi de la cartographie au backend FastAPI...")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("[-] Erreur de sérialisation JSON : %v", err)
	}

	// 5. Transmission HTTP Post au backend FastAPI pour injection Neo4j
	resp, err := http.Post("http://localhost:8000/scan", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("[-] Impossible de joindre le backend FastAPI sur http://localhost:8000 : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		fmt.Println("[+] Collecte réelle transmise et injectée avec succès dans Neo4j !")
	} else {
		fmt.Printf("[-] Échec de l'injection. Le backend a retourné le code : %d\n", resp.StatusCode)
	}
}
