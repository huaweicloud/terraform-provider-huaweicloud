package addons

import (
	"github.com/chnsz/golangsdk"
)

type ListAddon struct {
	// API type, fixed value "List"
	Kind string `json:"kind"`
	// API version, fixed value "v3"
	Apiversion string `json:"apiVersion"`
	// all Node Pools
	Addons []Addon `json:"items"`
}

type Addon struct {
	// API type, fixed value Addon
	Kind string `json:"kind"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion"`
	// Metadata of an Addon
	Metadata MetaData `json:"metadata"`
	// Specifications of an Addon
	Spec Spec `json:"spec"`
	// Status of an Addon
	Status Status `json:"status"`
}

//Metadata required to create an addon
type MetaData struct {
	// Addon unique name
	Name string `json:"name"`
	// Addon unique Id
	Id string `json:"uid"`
	// Addon tag, key/value pair format
	Labels map[string]string `json:"lables"`
	// Addon annotation, key/value pair format
	Annotations map[string]string `json:"annotaions"`
}

//Specifications to create an addon
type Spec struct {
	// For the addon version.
	Version string `json:"version"`
	// Cluster ID.
	ClusterID string `json:"clusterID"`
	// Addon Template Name.
	AddonTemplateName string `json:"addonTemplateName"`
	// Addon Template Type.
	AddonTemplateType string `json:"addonTemplateType"`
	// Addon Template Labels.
	AddonTemplateLables []string `json:"addonTemplateLables"`
	// Addon Description.
	Description string `json:"description"`
	// Addon Parameters
	Values Values `json:"values"`
}

type Status struct {
	//The state of the addon
	Status string `json:"status"`
	//Reasons for the addon to become current
	Reason string `json:"reason"`
	//Error Message
	Message string `json:"message"`
	//The target versions of the addon
	TargetVersions []string `json:"targetVersions"`
	//Current version of the addon
	CurrentVersion Versions `json:"currentVersion"`
}

type Versions struct {
	// Version of the addon
	Version string `json:"version"`
	// The installing param of the addon
	Input map[string]interface{} `json:"input"`
	// Wether it is a stable version
	Stable bool `json:"stable"`
	// Translate information
	Translate map[string]interface{} `json:"translate"`
	// Supported versions
	SupportVersions []SupportVersions `json:"supportVersions"`
}

type SupportVersions struct {
	ClusterType    string   `json:"clusterType"`
	ClusterVersion []string `json:"clusterVersion"`
}
type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts an Addon.
func (r commonResult) Extract() (*Addon, error) {
	var s Addon
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractAddon is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractAddon() ([]Addon, error) {
	var s ListAddon
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Addons, nil
}

// ListResult represents the result of a list operation. Call its ExtractAddon
// method to interpret it as a Addon.
type ListResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as an Addon.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as an Addon.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as an Addon.
type UpdataResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
