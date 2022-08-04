package responses

type Credentials struct {
	Results []WalletCredential `json:"results"`
}

type WalletCredential struct {
	Attrs     map[string]string `json:"attrs"`
	CredDefID string            `json:"cred_def_id"`
	CredRevID interface{}       `json:"cred_rev_id"`
	Referent  string            `json:"referent"`
	RevRegID  interface{}       `json:"rev_reg_id"`
	SchemaID  string            `json:"schema_id"`
}
