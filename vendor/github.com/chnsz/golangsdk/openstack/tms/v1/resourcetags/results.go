package resourcetags

// FailedResource is the structure that represents the resource list that set tags failed.
type FailedResource struct {
	// Specifies the resource ID.
	ResourceId string `json:"resource_id"`
	// Specifies the resource type.
	ResourceType string `json:"resource_type"`
	// Specifies the error code.
	ErrorCode string `json:"error_code"`
	// Specifies the error message.
	ErrorMsg string `json:"error_msg"`
}
