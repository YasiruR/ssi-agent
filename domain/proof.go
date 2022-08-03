package domain

type PresentationRequest struct {
	Indy struct {
		Name                string               `json:"name"`
		RequestedAttributes map[string]Attribute `json:"requested_attributes"`
		RequestedPredicates map[string]Predicate `json:"requested_predicates"`
		Version             string               `json:"version"`
	} `json:"indy"`
}

type Attribute struct {
	Name         string `json:"name"`
	Restrictions []struct {
		CredDefID string `json:"cred_def_id"`
	} `json:"restrictions"`
}

type Predicate struct {
	Name         string `json:"name"`
	PType        string `json:"p_type"`
	PValue       int64  `json:"p_value"`
	Restrictions []struct {
		CredDefID string `json:"cred_def_id"`
	} `json:"restrictions"`
}
