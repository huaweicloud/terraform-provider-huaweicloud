package elasticresourcepool

// AssociateQueueResp is the structure that represents response of the AssociateElasticResourcePool method.
type AssociateQueueResp struct {
	// Whether the request is successfully sent. Value true indicates that the request is successfully sent.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
}
