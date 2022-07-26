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

type AcceptInvitation struct {
	Accept              string `json:"accept"`
	Alias               string `json:"alias"`
	ConnectionID        string `json:"connection_id"`
	ConnectionProtocol  string `json:"connection_protocol"`
	CreatedAt           string `json:"created_at"`
	ErrorMsg            string `json:"error_msg"`
	InboundConnectionID string `json:"inbound_connection_id"`
	InvitationKey       string `json:"invitation_key"`
	InvitationMode      string `json:"invitation_mode"`
	InvitationMsgID     string `json:"invitation_msg_id"`
	MyDid               string `json:"my_did"`
	RequestID           string `json:"request_id"`
	Rfc23State          string `json:"rfc23_state"`
	RoutingState        string `json:"routing_state"`
	State               string `json:"state"`
	TheirDid            string `json:"their_did"`
	TheirLabel          string `json:"their_label"`
	TheirPublicDid      string `json:"their_public_did"`
	TheirRole           string `json:"their_role"`
	UpdatedAt           string `json:"updated_at"`
}

type Error struct {
	ID          string `json:"@id"`
	Type        string `json:"@type"`
	Description struct {
		Code string `json:"code"`
		En   string `json:"en"`
	} `json:"description"`
	Impact string `json:"impact"`
	Thread struct {
		Pthid string `json:"pthid"`
	} `json:"~thread"`
}
