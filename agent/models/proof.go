package models

import "github.com/YasiruR/agent/domain"

// ProofPresentation is used as the value type for storing presentation requests in memory map
type ProofPresentation struct {
	PresExID string
	PresReq  domain.PresentationRequest
}
