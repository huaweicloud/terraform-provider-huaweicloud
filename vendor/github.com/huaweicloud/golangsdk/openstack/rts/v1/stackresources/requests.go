package stackresources

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToStackResourceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the rts attributes you want to see returned.
type ListOpts struct {
	//Specifies the logical resource ID of the resource.
	LogicalID string `q:"logical_resource_id"`

	//Name is the human readable name for the Resource.
	Name string `q:"resource_name"`

	//Specifies the Physical resource ID of the resource.
	PhysicalID string `q:"physical_resource_id"`

	//Status indicates whether or not a subnet is currently operational.
	Status string `q:"resource_status"`

	//Specifies the resource type that are defined in the template.
	Type string `q:"resource_type"`
}

// List returns collection of
// resources. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those resources that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(client *golangsdk.ServiceClient, stackName string, opts ListOpts) ([]Resource, error) {
	u := listURL(client, stackName)
	pages, err := pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return ResourcePage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allResources, err := ExtractResources(pages)
	if err != nil {
		return nil, err
	}

	return FilterResources(allResources, opts)
}

func FilterResources(resources []Resource, opts ListOpts) ([]Resource, error) {

	var refinedResources []Resource
	var matched bool
	m := map[string]interface{}{}

	if opts.LogicalID != "" {
		m["LogicalID"] = opts.LogicalID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.PhysicalID != "" {
		m["PhysicalID"] = opts.PhysicalID
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.Type != "" {
		m["Type"] = opts.Type
	}

	if len(m) > 0 && len(resources) > 0 {
		for _, resource := range resources {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&resource, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedResources = append(refinedResources, resource)
			}
		}

	} else {
		refinedResources = resources
	}

	return refinedResources, nil
}

func getStructField(v *Resource, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
