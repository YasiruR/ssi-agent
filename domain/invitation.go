package domain

type Invitation struct {
	ID              string   `json:"@id"`
	Type            string   `json:"@type"`
	Label           string   `json:"label"`
	RecipientKeys   []string `json:"recipientKeys"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}
