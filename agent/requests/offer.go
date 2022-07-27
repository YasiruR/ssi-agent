package requests

import "github.com/YasiruR/agent/domain"

type Offer struct {
	//AutoIssue         bool                     `json:"auto_issue"`
	//AutoRemove        bool                     `json:"auto_remove"`
	Comment           string                   `json:"comment"`
	ConnectionID      string                   `json:"connection_id"`
	CredentialPreview domain.CredentialPreview `json:"credential_preview"`
	Filter            struct {
		Indy struct{} `json:"indy"`
	} `json:"filter"`

	//Filter            struct {
	//	Indy struct {
	//		CredDefID       string `json:"cred_def_id"`
	//		IssuerDid       string `json:"issuer_did"`
	//		SchemaID        string `json:"schema_id"`
	//		SchemaIssuerDid string `json:"schema_issuer_did"`
	//		SchemaName      string `json:"schema_name"`
	//		SchemaVersion   string `json:"schema_version"`
	//	} `json:"indy"`
	//	LdProof struct {
	//		Credential struct {
	//			Context           []string    `json:"@context"`
	//			CredentialSubject interface{} `json:"credentialSubject"`
	//			Description       string      `json:"description"`
	//			Identifier        string      `json:"identifier"`
	//			IssuanceDate      string      `json:"issuanceDate"`
	//			Issuer            string      `json:"issuer"`
	//			Name              string      `json:"name"`
	//			Type              []string    `json:"type"`
	//		} `json:"credential"`
	//		Options struct {
	//			ProofType string `json:"proofType"`
	//		} `json:"options"`
	//	} `json:"ld_proof"`
	//} `json:"filter"`
	//Trace bool `json:"trace"`
}
