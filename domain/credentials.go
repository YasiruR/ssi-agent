package domain

type CredentialPreview struct {
	Type       string `json:"@type"`
	Attributes []struct {
		Mime_type string `json:"mime-type"`
		Name      string `json:"name"`
		Value     string `json:"value"`
	} `json:"attributes"`
}
