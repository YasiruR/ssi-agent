package requests

type IssueCredentials struct {
	AutoIssue    bool   `json:"auto_issue"`
	AutoOffer    bool   `json:"auto_offer"`
	AutoRemove   bool   `json:"auto_remove"`
	ConnID       string `json:"conn_id"`
	CreatedAt    string `json:"created_at"`
	CredExID     string `json:"cred_ex_id"`
	CredIDStored string `json:"cred_id_stored"`
	CredIssue    struct {
		ID                 string `json:"@id"`
		Type               string `json:"@type"`
		Comment            string `json:"comment"`
		Credentials_attach []struct {
			ID   string `json:"@id"`
			Data struct {
				Base64 string `json:"base64"`
			} `json:"data"`
			Mime_type string `json:"mime-type"`
		} `json:"credentials~attach"`
		Formats []struct {
			AttachID string `json:"attach_id"`
			Format   string `json:"format"`
		} `json:"formats"`
		Thread struct {
			Thid string `json:"thid"`
		} `json:"~thread"`
	} `json:"cred_issue"`
	CredOffer struct {
		ID                string `json:"@id"`
		Type              string `json:"@type"`
		Comment           string `json:"comment"`
		CredentialPreview struct {
			Type       string `json:"@type"`
			Attributes []struct {
				Mime_type string `json:"mime-type"`
				Name      string `json:"name"`
				Value     string `json:"value"`
			} `json:"attributes"`
		} `json:"credential_preview"`
		Formats []struct {
			AttachID string `json:"attach_id"`
			Format   string `json:"format"`
		} `json:"formats"`
		Offers_attach []struct {
			ID   string `json:"@id"`
			Data struct {
				Base64 string `json:"base64"`
			} `json:"data"`
			Mime_type string `json:"mime-type"`
		} `json:"offers~attach"`
		Thread struct{} `json:"~thread"`
	} `json:"cred_offer"`
	CredPreview struct {
		Type       string `json:"@type"`
		Attributes []struct {
			Mime_type string `json:"mime-type"`
			Name      string `json:"name"`
			Value     string `json:"value"`
		} `json:"attributes"`
	} `json:"cred_preview"`
	CredProposal struct {
		ID                string `json:"@id"`
		Type              string `json:"@type"`
		Comment           string `json:"comment"`
		CredentialPreview struct {
			Type       string `json:"@type"`
			Attributes []struct {
				Mime_type string `json:"mime-type"`
				Name      string `json:"name"`
				Value     string `json:"value"`
			} `json:"attributes"`
		} `json:"credential_preview"`
		Filters_attach []struct {
			ID   string `json:"@id"`
			Data struct {
				Base64 string `json:"base64"`
			} `json:"data"`
			Mime_type string `json:"mime-type"`
		} `json:"filters~attach"`
		Formats []struct {
			AttachID string `json:"attach_id"`
			Format   string `json:"format"`
		} `json:"formats"`
	} `json:"cred_proposal"`
	Initiator string `json:"initiator"`
	Role      string `json:"role"`
	State     string `json:"state"`
	ThreadID  string `json:"thread_id"`
	Trace     bool   `json:"trace"`
	UpdatedAt string `json:"updated_at"`
}

type IssueCredentialsIndy struct {
	CreatedAt           string `json:"created_at"`
	CredExID            string `json:"cred_ex_id"`
	CredExIndyID        string `json:"cred_ex_indy_id"`
	CredRequestMetadata struct {
		MasterSecretBlindingData struct {
			VPrime  string      `json:"v_prime"`
			VrPrime interface{} `json:"vr_prime"`
		} `json:"master_secret_blinding_data"`
		MasterSecretName string `json:"master_secret_name"`
		Nonce            string `json:"nonce"`
	} `json:"cred_request_metadata"`
	RevRegID  string `json:"rev_reg_id"`
	UpdatedAt string `json:"updated_at"`
}
