package requests

import "github.com/YasiruR/agent/domain"

type Offer struct {
	//AutoIssue         bool                     `json:"auto_issue"`
	AutoRemove        bool                     `json:"auto_remove"`
	Comment           string                   `json:"comment"`
	ConnectionID      string                   `json:"connection_id"`
	CredentialPreview domain.CredentialPreview `json:"credential_preview"`
	Filter            struct {
		Indy domain.IndySchemaMeta `json:"indy"`
	} `json:"filter"`
}
