package requests

import "github.com/YasiruR/agent/domain"

type PresentationProof struct {
	AutoPresent bool `json:"auto_present"`
	ByFormat    struct {
		PresRequest domain.PresentationRequest `json:"pres_request"`
	} `json:"by_format"`
	ConnectionID string `json:"connection_id"`
	CreatedAt    string `json:"created_at"`
	Initiator    string `json:"initiator"`
	PresExID     string `json:"pres_ex_id"`
	PresRequest  struct {
		ID      string `json:"@id"`
		Type    string `json:"@type"`
		Comment string `json:"comment"`
		Formats []struct {
			AttachID string `json:"attach_id"`
			Format   string `json:"format"`
		} `json:"formats"`
		RequestPresentations_attach []struct {
			ID   string `json:"@id"`
			Data struct {
				Base64 string `json:"base64"`
			} `json:"data"`
			Mime_type string `json:"mime-type"`
		} `json:"request_presentations~attach"`
		WillConfirm bool `json:"will_confirm"`
	} `json:"pres_request"`
	Role      string `json:"role"`
	State     string `json:"state"`
	ThreadID  string `json:"thread_id"`
	Trace     bool   `json:"trace"`
	UpdatedAt string `json:"updated_at"`
}
