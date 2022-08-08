package requests

type CreateInvitation struct {
	Alias       string `json:"alias"`
	Attachments []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"attachments"`
	HandshakeProtocols []string `json:"handshake_protocols"`
	MediationID        string   `json:"mediation_id"`
	Metadata           struct{} `json:"metadata"`
	MyLabel            string   `json:"my_label"`
	UsePublicDid       bool     `json:"use_public_did"`
}
