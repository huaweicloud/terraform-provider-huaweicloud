package metadatas

type MetadataResult struct {
	Message string `json:"message"`
}

type Metadata struct {
	ID           string `json:"id"`
	IdpID        string `json:"idp_id"`
	EntityID     string `json:"entity_id"`
	ProtocolID   string `json:"protocol_id"`
	DomainID     string `json:"domain_id"`
	XAccountType string `json:"xaccount_type"`
	UpdateTime   string `json:"update_time"`
	Data         string `json:"data"`
}
