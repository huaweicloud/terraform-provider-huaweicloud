package flavors

import (
	"github.com/huaweicloud/golangsdk/pagination"
)

// FlavorPage is a single page of Flavor results.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the page contains no results.
func (r FlavorPage) IsEmpty() (bool, error) {
	s, err := ExtractFlavors(r)
	return len(s) == 0, err
}

// ExtractFlavors extracts a slice of Flavors from a List result.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorPage)).ExtractInto(&s)
	return s.Flavors, err
}

// Flavor represents a DDS flavor.
type Flavor struct {
	EngineName       string   `json:"engine_name"`
	Type             string   `json:"type"`
	Vcpus            string   `json:"vcpus"`
	Ram              string   `json:"ram"`
	SpecCode         string   `json:"spec_code"`
	AvailabilityZone []string `json:"availability_zone"`
}
