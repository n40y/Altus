// scanner/collection/aws.go
package collection

import "altus/scanner/models"

// DiscoverResources simule la découverte active de l'infrastructure Cloud
func DiscoverResources() models.ScanResult {
	var result models.ScanResult

	// Découverte des Nœuds
	result.Nodes = append(result.Nodes, models.Node{Name: "go-scanner-user", Type: "IAM_USER"})
	result.Nodes = append(result.Nodes, models.Node{Name: "AssumeRolePolicy", Type: "IAM_POLICY"})
	result.Nodes = append(result.Nodes, models.Node{Name: "altus-crypto-vault", Type: "S3_BUCKET"})
	result.Nodes = append(result.Nodes, models.Node{Name: "ec2-api-gateway", Type: "EC2_INSTANCE"})

	// Découverte des Relations
	result.Relations = append(result.Relations, models.Relation{
		Source: "go-scanner-user", Target: "AssumeRolePolicy", Relation: "ALLOW_ASSUME",
	})
	result.Relations = append(result.Relations, models.Relation{
		Source: "AssumeRolePolicy", Target: "altus-crypto-vault", Relation: "READ_WRITE",
	})
	result.Relations = append(result.Relations, models.Relation{
		Source: "ec2-api-gateway", Target: "AssumeRolePolicy", Relation: "USES_ROLE",
	})

	return result
}
