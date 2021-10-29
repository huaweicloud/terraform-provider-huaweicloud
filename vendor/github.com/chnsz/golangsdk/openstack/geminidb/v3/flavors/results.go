package flavors

import (
	"github.com/chnsz/golangsdk/pagination"
)

type Flavor struct {
	EngineName       string   `json:"engine_name"`
	EngineVersion    string   `json:"engine_version"`
	Vcpus            string   `json:"vcpus"`
	Ram              string   `json:"ram"`
	SpecCode         string   `json:"spec_code"`
	AvailabilityZone []string `json:"availability_zone"`
	// AZ status
	AzStatus map[string]string `json:"az_status"`
}

type ListFlavorResponse struct {
	Flavors    []Flavor `json:"flavors"`
	TotalCount int      `json:"total_count"`
}

type FlavorPage struct {
	pagination.SinglePageBase
}

func (r FlavorPage) IsEmpty() (bool, error) {
	data, err := ExtractFlavors(r)
	if err != nil {
		return false, err
	}
	return len(data.Flavors) == 0, err
}

func ExtractFlavors(r pagination.Page) (ListFlavorResponse, error) {
	var s ListFlavorResponse
	err := (r.(FlavorPage)).ExtractInto(&s)
	return s, err
}
