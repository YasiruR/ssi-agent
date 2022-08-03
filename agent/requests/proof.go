package requests

import "github.com/YasiruR/agent/domain"

type ProofRequest struct {
	Comment      string                     `json:"comment"`
	ConnectionID string                     `json:"connection_id"`
	PresentReq   domain.PresentationRequest `json:"presentation_request"`
}
