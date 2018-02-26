package backendmember

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Backend is the primary load balancing configuration object that specifies
// the loadbalancer and port on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type Backend struct {
	// Specifies the private IP address of the backend ECS.
	ServerAddress string `json:"server_address"`
	// Specifies the backend ECS ID.
	ID string `json:"id"`
	// Specifies the floating IP address assigned to the backend ECS.
	Address string `json:"address"`
	// Specifies the backend ECS status. The value is ACTIVE, PENDING, or ERROR.
	Status string `json:"status"`
	// Specifies the health check status. The value is NORMAL, ABNORMAL, or UNAVAILABLE.
	HealthStatus string `json:"health_status"`
	// Specifies the time when information about the backend ECS was updated.
	UpdateTime string `json:"update_time"`
	// Specifies the time when the backend ECS was created.
	CreateTime string `json:"create_time"`
	// Specifies the backend ECS name.
	ServerName string `json:"server_name"`
	// Specifies the original back member ID.
	ServerID string `json:"server_id"`
	// Specifies the listener to which the backend ECS belongs.
	Listeners []map[string]interface{} `json:"listeners"`
}

// ListenerPage is the page returned by a pager when traversing over a
// collection of routers.
type BackendPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackendPage) NextPageURL() (string, error) {
	return "", nil
}

// IsEmpty checks whether a RouterPage struct is empty.
func (r BackendPage) IsEmpty() (bool, error) {
	is, err := ExtractBackend(r)
	return len(is) == 0, err
}

// ExtractBackend accepts a Page struct, specifically a ListenerPage struct,
// and extracts the elements into a slice of Listener structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackend(r pagination.Page) ([]Backend, error) {
	var Backends []Backend
	err := (r.(BackendPage)).ExtractInto(&Backends)
	return Backends, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) Extract() (*Backend, error) {
	//fmt.Printf("Extracting Backend...\n")
	var Backends []Backend
	err := r.ExtractInto(&Backends)
	if err != nil {
		//fmt.Printf("Error: %s.\n", err.Error())
		return nil, err
	} else {
		if len(Backends) == 0 {
			return nil, nil
		}
		return &Backends[0], nil
	}
}

// AddResult represents the result of a create operation.
type AddResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

func (r GetResult) Extract() ([]Backend, error) {
	var b []Backend
	err := r.ExtractInto(&b)
	return b, err
}

// RemoveResult represents the result of a delete operation.
type RemoveResult struct {
	golangsdk.ErrResult
}
