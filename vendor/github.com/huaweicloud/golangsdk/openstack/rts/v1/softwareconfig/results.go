package softwareconfig

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type SoftwareConfig struct {
	// Specifies the software configuration input.
	Inputs []map[string]interface{} `json:"inputs"`
	//Specifies the name of the software configuration.
	Name string `json:"name"`
	//Specifies the software configuration output.
	Outputs []map[string]interface{} `json:"outputs"`
	//Specifies the time when a configuration is created.
	CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
	//Specifies the name of the software configuration group.
	Group string `json:"group"`
	//Specifies the configuration code.
	Config string `json:"config"`
	//Specifies configuration options.
	Options map[string]interface{} `json:"options"`
	//Specifies the software configuration ID.
	Id string `json:"id"`
}

// SoftwareConfigPage is the page returned by a pager when traversing over a
// collection of Software Configurations.
type SoftwareConfigPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of Software Configs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SoftwareConfigPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"software_config_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SoftwareConfigPage struct is empty.
func (r SoftwareConfigPage) IsEmpty() (bool, error) {
	is, err := ExtractSoftwareConfigs(r)
	return len(is) == 0, err
}

// ExtractSoftwareConfigs accepts a Page struct, specifically a SoftwareConfigPage struct,
// and extracts the elements into a slice of Software Configs structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSoftwareConfigs(r pagination.Page) ([]SoftwareConfig, error) {
	var s struct {
		SoftwareConfigs []SoftwareConfig `json:"software_configs"`
	}
	err := (r.(SoftwareConfigPage)).ExtractInto(&s)
	return s.SoftwareConfigs, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a Software configuration.
func (r commonResult) Extract() (*SoftwareConfig, error) {
	var s struct {
		SoftwareConfig *SoftwareConfig `json:"software_config"`
	}
	err := r.ExtractInto(&s)
	return s.SoftwareConfig, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Software configuration.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Software configuration.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Software configuration.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
