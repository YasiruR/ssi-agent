package requests

import "github.com/YasiruR/agent/domain"

type ProofRequest struct {
	Comment      string                     `json:"comment"`
	ConnectionID string                     `json:"connection_id"`
	PresentReq   domain.PresentationRequest `json:"presentation_request"`
}

type ProofPresentation struct {
	Indy struct {
		RequestedAttributes    map[string]AdditionalProp `json:"requested_attributes"`
		RequestedPredicates    map[string]AdditionalProp `json:"requested_predicates"`
		SelfAttestedAttributes map[string]string         `json:"self_attested_attributes"`
		Trace                  bool                      `json:"trace"`
	} `json:"indy"`
}

type AdditionalProp struct {
	CredID   string `json:"cred_id"`
	Revealed bool   `json:"revealed"`
}
