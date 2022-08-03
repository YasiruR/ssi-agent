package requests

import "github.com/YasiruR/agent/domain"

type Offer struct {
	AutoProcess bool                     `json:"auto_process"`
	CredPreview domain.CredentialPreview `json:"credential_preview"`
	Filter      struct {
		Indy domain.IndySchemaMeta `json:"indy"`
	} `json:"filter"`
}
