package instances

// CreateResp is the structure that represents the response of the API request.
type CreateResp struct {
	// Microservice instance ID.
	ID string `json:"instanceId"`
}

// Instance is the structure that represents the details of the Microservice service instance.
type Instance struct {
	// Host information.
	HostName string `json:"hostName"`
	// Access address information.
	Endpoints []string `json:"endpoints"`
	// Instance ID, which must be unique. The instance ID is generated by the service center.
	ID string `json:"instanceId"`
	// Microservice ID, which must be unique.
	// During instance creation, use the service ID in the URL instead of the service ID here.
	ServiceId string `json:"serviceId"`
	// Microservice version.
	Version string `json:"version"`
	// Instance status. Value: UP, DOWN, STARTING, or OUTOFSERVICE. Default value: UP.
	Status string `json:"status"`
	// Extended attribute. You can customize a key and value. The value must be at least 1 byte long.
	Properties map[string]interface{} `json:"properties"`
	// Health check information.
	HealthCheck HealthCheck `json:"healthCheck"`
	// Data center information.
	DataCenterInfo DataCenter `json:"dataCenterInfo"`
	// Time when an instance is created, which is automatically generated.
	Timestamp string `json:"timestamp"`
	// Update time.
	ModTimestamp string `json:"modTimestamp"`
}

// ErrorResponse is the structure that represents the details of the request error.
type ErrorResponse struct {
	// Error detail.
	Detail string `json:"detail"`
	// Error code.
	ErrCode string `json:"errorCode"`
	// Error message.
	ErrMessage string `json:"errorMessage"`
}
