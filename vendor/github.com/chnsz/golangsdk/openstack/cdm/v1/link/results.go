package link

type LinkCreateResponse struct {
	Name             string             `json:"name"`
	ValidationResult []validationResult `json:"validation-result"`
}

type validationResult struct {
	LinkConfig []LinkValidationDetail `json:"linkConfig"`
}

type LinkValidationDetail struct {
	Message string `json:"message"`
	// ERROR,WARNING
	Status string `json:"status"`
}

type LinkDetail struct {
	// Link list. For details, see the description of the links parameter.
	Links []Link `json:"links"`
	// Source and destination data sources not supported by table/file migration
	FromToUnMapping string `json:"fromTo-unMapping"`
	// Source and destination data sources supported by entire DB migration
	BatchFromToMapping string `json:"batchFromTo-mapping"`
}

type LinkUpdateResponse struct {
	ValidationResult []validationResult `json:"validation-result"`
}

type LinkDeleteResponse struct {
	// Error code
	ErrCode string `json:"errCode"`
	// Error message
	ErrMessage string `json:"externalMessage"`
}
