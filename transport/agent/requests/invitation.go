package requests

import "github.com/YasiruR/agent/domain"

type AcceptInvitation struct {
	ConnectionID  string            `json:"connection_id"`
	Invitation    domain.Invitation `json:"invitation"`
	InvitationURL string            `json:"invitation_url"`
}

//type AcceptInvitation struct {
//	CreatedAt     string            `json:"created_at"`
//	InviMsgID     string            `json:"invi_msg_id"`
//	Invitation    domain.Invitation `json:"invitation"`
//	InvitationID  string            `json:"invitation_id"`
//	InvitationURL string            `json:"invitation_url"`
//	OobID         string            `json:"oob_id"`
//	State         string            `json:"state"`
//	Trace         bool              `json:"trace"`
//	UpdatedAt     string            `json:"updated_at"`
//}
