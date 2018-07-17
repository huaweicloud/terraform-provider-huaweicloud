package stackresources

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Resource represents a stack resource.
type Resource struct {
	CreationTime time.Time        `json:"-"`
	Links        []golangsdk.Link `json:"links"`
	LogicalID    string           `json:"logical_resource_id"`
	Name         string           `json:"resource_name"`
	PhysicalID   string           `json:"physical_resource_id"`
	RequiredBy   []string         `json:"required_by"`
	Status       string           `json:"resource_status"`
	StatusReason string           `json:"resource_status_reason"`
	Type         string           `json:"resource_type"`
	UpdatedTime  time.Time        `json:"-"`
}

// ResourcePage is the page returned by a pager when traversing over a
// collection of resources.
type ResourcePage struct {
	pagination.LinkedPageBase
}

func (r *Resource) UnmarshalJSON(b []byte) error {
	type tmp Resource
	var s struct {
		tmp
		CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
		UpdatedTime  golangsdk.JSONRFC3339NoZ `json:"updated_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Resource(s.tmp)

	r.CreationTime = time.Time(s.CreationTime)
	r.UpdatedTime = time.Time(s.UpdatedTime)

	return nil
}

// IsEmpty returns true if a page contains no Server results.
func (r ResourcePage) IsEmpty() (bool, error) {
	resources, err := ExtractResources(r)
	return len(resources) == 0, err
}

// ExtractResources accepts a Page struct, specifically a ResourcePage struct,
// and extracts the elements into a slice of Resource structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractResources(r pagination.Page) ([]Resource, error) {
	var s struct {
		Resources []Resource `json:"resources"`
	}
	err := (r.(ResourcePage)).ExtractInto(&s)
	return s.Resources, err
}
