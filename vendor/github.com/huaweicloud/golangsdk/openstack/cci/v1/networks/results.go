package networks

import (
	"github.com/huaweicloud/golangsdk"
)

type Network struct {
	//API type, fixed value Network
	Kind string `json:"kind"`
	//API version, fixed value networking.cci.io
	ApiVersion string `json:"apiVersion"`
	//Metadata of a Network
	Metadata MetaData `json:"metadata"`
	//Specifications of a Network
	Spec Spec `json:"spec"`
	//Status of a Network
	Status Status `json:"status"`
}

//Metadata required to create a network
type MetaData struct {
	//Network unique name
	Name string `json:"name"`
	//Network unique Id
	Id string `json:"uid"`
	//Network annotation, key/value pair format
	Annotations map[string]string `json:"annotations"`
}

type Status struct {
	//The state of the network
	State string `json:"state"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a network.
func (r commonResult) Extract() (*Network, error) {
	var s Network
	err := r.ExtractInto(&s)
	return &s, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Network.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Network.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
