package responses

type Invitation struct {
	CreatedAt  string `json:"created_at"`
	InviMsgID  string `json:"invi_msg_id"`
	Invitation struct {
		ID                 string   `json:"@id"`
		Type               string   `json:"@type"`
		HandshakeProtocols []string `json:"handshake_protocols"`
		Label              string   `json:"label"`
		Requests_attach    []struct {
			ID        string `json:"@id"`
			ByteCount int64  `json:"byte_count"`
			Data      struct {
				Base64 string `json:"base64"`
				JSON   struct {
					Sample string `json:"sample"`
				} `json:"json"`
				Jws struct {
					Header struct {
						Kid string `json:"kid"`
					} `json:"header"`
					Protected  string `json:"protected"`
					Signature  string `json:"signature"`
					Signatures []struct {
						Header struct {
							Kid string `json:"kid"`
						} `json:"header"`
						Protected string `json:"protected"`
						Signature string `json:"signature"`
					} `json:"signatures"`
				} `json:"jws"`
				Links  []string `json:"links"`
				Sha256 string   `json:"sha256"`
			} `json:"data"`
			Description string `json:"description"`
			Filename    string `json:"filename"`
			LastmodTime string `json:"lastmod_time"`
			Mime_type   string `json:"mime-type"`
		} `json:"requests~attach"`
		Services []struct {
			Did             string   `json:"did"`
			ID              string   `json:"id"`
			RecipientKeys   []string `json:"recipientKeys"`
			RoutingKeys     []string `json:"routingKeys"`
			ServiceEndpoint string   `json:"serviceEndpoint"`
			Type            string   `json:"type"`
		} `json:"services"`
	} `json:"invitation"`
	InvitationID  string `json:"invitation_id"`
	InvitationURL string `json:"invitation_url"`
	OobID         string `json:"oob_id"`
	State         string `json:"state"`
	Trace         bool   `json:"trace"`
	UpdatedAt     string `json:"updated_at"`
}
