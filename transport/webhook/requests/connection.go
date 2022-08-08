package requests

type Connections struct {
	Accept         string `json:"accept"`
	Alias          string `json:"alias"`
	ConnectionID   string `json:"connection_id"`
	CreatedAt      string `json:"created_at"`
	InvitationMode string `json:"invitation_mode"`
	Rfc23State     string `json:"rfc23_state"`
	RoutingState   string `json:"routing_state"`
	State          string `json:"state"`
	TheirLabel     string `json:"their_label"`
	TheirRole      string `json:"their_role"`
	UpdatedAt      string `json:"updated_at"`
}
