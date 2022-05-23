package services

// CreateResp is the structure that represents the response of the API request.
type CreateResp struct {
	// Microservice ID.
	ID string `json:"serviceId"`
}

// ServiceResp is an object that specifies the microservice configuration details.
type ServiceResp struct {
	// Microservice name, which must be unique in an application. The value contains 1 to 128 characters.
	// Regular expression: ^[a-zA-Z0-9]*$|^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$
	Name string `json:"serviceName" required:"true"`
	// Application ID, which must be unique. The value contains 1 to 160 characters.
	// Regular expression: ^[a-zA-Z0-9]*$|^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$
	AppId string `json:"appId" required:"true"`
	// Microservice version. The value contains 1 to 64 characters. Regular expression: ^[0-9]$|^[0-9]+(.[0-9]+)$
	Version string `json:"version" required:"true"`
	// Microservice ID, which must be unique. The value contains 1 to 64 characters. Regular expression: ^.*$
	ID string `json:"serviceId,omitempty"`
	// Service stage. Value: development, testing, acceptance, or production.
	// Only when the value is development, testing, or acceptance, you can use the API for uploading schemas in batches
	// to add or modify an existing schema. Default value: development.
	Environment string `json:"environment,omitempty"`
	// Microservice description. The value contains a maximum of 256 characters.
	Description string `json:"description,omitempty"`
	// Microservice level. Value: FRONT, MIDDLE, or BACK.
	Level string `json:"level,omitempty"`
	// Microservice registration mode. Value: SDK, PLATFORM, SIDECAR, or UNKNOWN.
	RegisterBy string `json:"registerBy,omitempty"`
	// Foreign key ID of a microservice access schema. The array length supports a maximum of 100 schemas.
	Schemas []string `json:"schemas,omitempty"`
	// Microservice status. Value: UP or DOWN. Default value: UP.
	Status string `json:"status,omitempty"`
	// Microservice registration time.
	Timestamp string `json:"timestamp,omitempty"`
	// Latest modification time (UTC).
	ModTimestamp string `json:"modTimestamp,omitempty"`
	// Development framework.
	Framework Framework `json:"framework,omitempty"`
	// Service path.
	Paths []ServicePath `json:"paths,omitempty"`
}
