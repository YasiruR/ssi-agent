package domain

type IndySchemaMeta struct {
	CredDefID       string `json:"cred_def_id"`
	IssuerDid       string `json:"issuer_did"`
	SchemaID        string `json:"schema_id"`
	SchemaIssuerDid string `json:"schema_issuer_did"`
	SchemaName      string `json:"schema_name"`
	SchemaVersion   string `json:"schema_version"`
}
