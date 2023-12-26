package connections

import "github.com/chnsz/golangsdk/pagination"

type createResp struct {
	// The ID of the data connection.
	DataConnectionId string `json:"data_connection_id"`
}

// Connection is the structure that represents the data connection detail.
type Connection struct {
	// The data connection name.
	DwName string `json:"dw_name"`
	// The data connection type.
	DwType string `json:"dw_type"`
	// The dynamic configuration for the specified type of data connection.
	DwConfig interface{} `json:"dw_config"`
	// The agent ID.
	AgentId string `json:"agent_id"`
	// The agent name.
	AgentName string `json:"agent_name"`
	// The data connection mode.
	EnvType int `json:"env_type"`
	// The qualified name of the data connection.
	QualifiedName string `json:"qualified_name"`
	// The data connection ID.
	DwId string `json:"dw_id"`
	// The creator name.
	CreateUser string `json:"create_user"`
	// The creation time of the data connection.
	CreateTime int `json:"create_time"`
	// The catagory of the data connection.
	DwCatagory string `json:"dw_catagory"`
	// The update type of the data connection.
	UpdateType int `json:"update_type"`
}

// ValidateResp is the structure that represents the result of data connection pre-check.
type ValidateResp struct {
	// The message of the data connection pre-check.
	Message string `json:"message"`
	// Whether the data connection pre-check is successful.
	IsSuccess bool `json:"is_success"`
}

// ConnectionPage represents the response pages of the List method.
type ConnectionPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ConnectionPage struct is empty.
func (b ConnectionPage) IsEmpty() (bool, error) {
	arr, err := extractConnections(b)
	return len(arr) == 0, err
}

// ExtractConnections is a method to extract the list of data connections.
func extractConnections(r pagination.Page) ([]Connection, error) {
	var s []Connection
	err := r.(ConnectionPage).Result.ExtractIntoSlicePtr(&s, "data_connection_lists")
	return s, err
}
