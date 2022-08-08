package requests

import "github.com/YasiruR/agent/domain"

type ProofReq struct {
	PresentReq domain.PresentationRequest `json:"presentation_request"`
}
