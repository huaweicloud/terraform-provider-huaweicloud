package protocols

type IdentityProtocol struct {
	ID        string        `json:"id"`
	MappingID string        `json:"mapping_id"`
	Links     ProtocolLinks `json:"links"`
}

type ProtocolLinks struct {
	IdentityProvider string `json:"identity_provider"`
	Self             string `json:"self"`
}
