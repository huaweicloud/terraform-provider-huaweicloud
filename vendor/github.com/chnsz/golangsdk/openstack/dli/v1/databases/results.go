package databases

// RequestResp is a object that represents the result of Create and UpdateDBOwner operation.
type RequestResp struct {
	// Whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
}

// ListResp is a object that represents the result of List operation.
type ListResp struct {
	// Indicates whether the request is successfully executed.
	// Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
	// Total number of databases.
	DatabaseCount int `json:"database_count"`
	// Database information.
	Databases []Database `json:"databases"`
}

// Database is a object that represents the database detail.
type Database struct {
	// Name of a database.
	Name string `json:"database_name"`
	// Creator of a database.
	Owner string `json:"owner"`
	// Number of tables in a database.
	TableNumber int `json:"table_number"`
	// Information about a database.
	Description string `json:"description"`
	// Whether database is shared.
	IsShared bool `json:"is_shared"`
	// Enterprise project ID. The value 0 indicates the default enterprise project.
	// NOTE: Users who have enabled Enterprise Management can set this parameter to bind a specified project.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Resource ID.
	ResourceId string `json:"resource_id"`
}
