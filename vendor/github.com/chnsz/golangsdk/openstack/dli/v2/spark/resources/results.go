package resources

// Group is a object that represents the group create or group list operations.
type Group struct {
	// Module name.
	GroupName string `json:"group_name"`
	// Status of a package group to be uploaded.
	Status string `json:"status"`
	// List of names of resource packages contained in the group.
	Resources []string `json:"resources"`
	// Details about a group resource package.
	Details []Detail `json:"details"`
	// UNIX timestamp when a package group is uploaded.
	CreateTime int `json:"create_time"`
	// UNIX timestamp when a package group is updated.
	UpdateTime int `json:"update_time"`
	// Whether to upload resource packages in asynchronous mode.
	// The default value is false, indicating that the asynchronous mode is not used.
	// You are advised to upload resource packages in asynchronous mode.
	IsAsync bool `json:"is_async"`
	// Owner of a resource package.
	Owner string `json:"owner"`
	// The description, moduleName and moduleType are supported by function response of specified file upload.
	// Description of a resource module.
	Description string `json:"description"`
	// Name of a resource module.
	ModuleName string `json:"module_name"`
	// Type of a resource module.
	//   jar: User JAR file
	//   pyFile: User Python file
	//   file: User file
	ModuleType string `json:"module_type"`
}

// Group is a object that represents the detail about a group resource package.
type Detail struct {
	// UNIX time when a resource package is uploaded. The timestamp is expressed in milliseconds.
	CreateTime int `json:"create_time" required:"true"`
	// UNIX time when the uploaded resource package is uploaded. The timestamp is expressed in milliseconds.
	UpdateTime int `json:"update_time"`
	// Resource type.
	ResourceType string `json:"resource_type" required:"true"`
	// Resource name.
	ResourceName string `json:"resource_name"`
	// Value UPLOADING indicates that the resource package group is being uploaded.
	// Value READY indicates that the resource package has been uploaded.
	// Value FAILED indicates that the resource package fails to be uploaded.
	Status string `json:"status"`
	// Name of the resource packages in a queue.
	UnderlyingName string `json:"underlying_name"`
	// Whether to upload resource packages in asynchronous mode.
	// The default value is false, indicating that the asynchronous mode is not used.
	// You are advised to upload resource packages in asynchronous mode.
	IsAsync bool `json:"is_async"`
}

// ListResp is a object that represents the List method result.
type ListResp struct {
	// List of names, type and other informations of uploaded user resources.
	Resources []Resource `json:"resources"`
	// List of built-in resource groups. For details about the groups, see Table 5.
	Modules []Module `json:"modules"`
	// Uploaded package groups of a user.
	Groups []Group `json:"groups"`
	// Total number of returned resource packages.
	Total int `json:"total" required:"true"`
}

// Resource is a object that represents the names, type and other informations of uploaded user resources.
type Resource struct {
	// UNIX timestamp when a resource package is uploaded.
	CreateTime int `json:"create_time"`
	// UNIX timestamp when the uploaded resource package is uploaded.
	UpdateTime int `json:"update_time"`
	// Resource type.
	ResourceType string `json:"resource_type"`
	// Resource name.
	ResourceName string `json:"resource_name"`
	// Value UPLOADING indicates that the resource package is being uploaded.
	// Value READY indicates that the resource package has been uploaded.
	// Value FAILED indicates that the resource package fails to be uploaded.
	Status string `json:"status"`
	// Name of the resource package in the queue.
	UnderlyingName string `json:"underlying_name"`
	// Owner of a resource package.
	Owner string `json:"owner"`
}

// Module is a object that represents the built-in resource group.
type Module struct {
	// Module name.
	ModuleName string `json:"module_name"`
	// Module type.
	ModuleType string `json:"module_type"`
	// Value UPLOADING indicates that the package group is being uploaded.
	// Value READY indicates that the package group has been uploaded.
	// Value FAILED indicates that the package group fails to be uploaded.
	Status string `json:"status"`
	// List of names of resource packages contained in the group.
	Resources []string `json:"resources"`
	// Module description.
	Description string `json:"description"`
	// UNIX timestamp when a package group is uploaded.
	CreateTime int `json:"create_time"`
	// UNIX timestamp when a package group is updated.
	UpdateTime int `json:"update_time"`
}

// UpdateResp is a object that represents the update method result.
type UpdateResp struct {
	// Whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
}
