// scanne/models/models.go
package models

type Node struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Relation struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Relation string `json:"relation"`
}

type ScanResult struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}
