package domain

type Invitation struct {
	ID              string   `json:"@id"`
	Type            string   `json:"@type"`
	Label           string   `json:"label"`
	RecipientKeys   []string `json:"recipientKeys"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}

//type Invitation struct {
//	Type     string    `json:"@type"`
//	ID       string    `json:"@id"`
//	Label    string    `json:"label"`
//	Services []Service `json:"services"`
//	//RequestsAttach     []struct{} `json:"requests~attach"`
//	HandshakeProtocols []string `json:"handshake_protocols"`
//}
//
//type Service struct {
//	//Did             string      `json:"did"`
//	ID            string   `json:"id"`
//	Type          string   `json:"type"`
//	RecipientKeys []string `json:"recipientKeys"`
//	//RoutingKeys     interface{} `json:"routingKeys"`
//	ServiceEndpoint string `json:"serviceEndpoint"`
//}
