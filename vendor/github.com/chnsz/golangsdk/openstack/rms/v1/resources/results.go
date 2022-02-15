package resources

import (
	"github.com/chnsz/golangsdk/pagination"
)

type ListResp struct {
	Resources []Resource `json:"resources"`
	PageInfo  PageInfo   `json:"page_info"`
}

type Resource struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Provider          string                 `json:"provider"`
	Type              string                 `json:"type"`
	RegionId          string                 `json:"region_id"`
	ProjectId         string                 `json:"project_id"`
	ProjectName       string                 `json:"project_name"`
	EpId              string                 `json:"ep_id"`
	EpName            string                 `json:"ep_name"`
	Checksum          string                 `json:"checksum"`
	Created           string                 `json:"created"`
	Updated           string                 `json:"updated"`
	ProvisioningState string                 `json:"provisioning_state"`
	Tags              map[string]string      `json:"tags"`
	Properties        map[string]interface{} `json:"properties"`
}

type PageInfo struct {
	CurrentCount int    `json:"current_count"`
	NextMarker   string `json:"next_marker"`
}

// ResourcePage is the page returned by a pager when traversing over a
// collection of route tables
type ResourcePage struct {
	pagination.MarkerPageBase
}

// LastMarker returns the last resource ID in a ListResult
func (r ResourcePage) LastMarker() (string, error) {
	var s ListResp
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.PageInfo.NextMarker, nil
}

// IsEmpty checks whether a ResourcePage struct is empty.
func (r ResourcePage) IsEmpty() (bool, error) {
	tables, err := ExtractResources(r)
	return len(tables) == 0, err
}

// ExtractResources accepts a Page struct, specifically a ResourcePage struct,
// and extracts the elements into a slice of Resource structs.
func ExtractResources(r pagination.Page) ([]Resource, error) {
	var s ListResp
	err := (r.(ResourcePage)).ExtractInto(&s)
	return s.Resources, err
}
