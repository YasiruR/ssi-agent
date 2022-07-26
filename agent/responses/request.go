package responses

type ConnRequest struct {
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
