// scanner/collection/security.go

package collection

import "altus/scanner/models"

// RunSecurityAudit exécute la logique de détection d'exposition publique
func RunSecurityAudit(data *models.ScanResult) {
	if len(data.Nodes) == 0 {
		return
	}

	hasInternet := false
	for _, n := range data.Nodes {
		if n.Name == "INTERNET" {
			hasInternet = true
			break
		}
	}
	if !hasInternet {
		data.Nodes = append(data.Nodes, models.Node{
			Name: "INTERNET",
			Type: "GATEWAY",
		})
	}

	// Lie automatiquement les S3 ou EC2 trouvés pour visualiser la surface d'exposition
	for _, node := range data.Nodes {
		if node.Type == "S3_BUCKET" || node.Type == "EC2_INSTANCE" {
			data.Relations = append(data.Relations, models.Relation{
				Source:   "INTERNET",
				Target:   node.Name,
				Relation: "EXPOSED_TO",
			})
		}
	}
}
