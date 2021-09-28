package flavors

import (
	"github.com/chnsz/golangsdk"
)

// Flavor for dcs
type Flavor struct {
	SpecCode         string           `json:"spec_code"`
	CacheMode        string           `json:"cache_mode"`
	Engine           string           `json:"engine"`
	EngineVersion    string           `json:"engine_version"`
	TenantIPCount    int              `json:"tenant_ip_count"`
	IsDes            bool             `json:"is_dec"`
	Capacity         []string         `json:"capacity"`
	BillingMode      []string         `json:"billing_mode"`
	ProductType      string           `json:"product_type"`
	CPUType          string           `json:"cpu_type"`
	StorageType      string           `json:"storage_type"`
	PricingType      string           `json:"pricing_type"`
	ServiceTypeCode  string           `json:"cloud_service_type_code"`
	ResourceTypeCode string           `json:"cloud_resource_type_code"`
	Attributes       []AttrsObject    `json:"attrs"`
	AvailableZones   []FlavorAzObject `json:"flavors_available_zones"`
}

// AttrsObject contains attributes of the flavor
type AttrsObject struct {
	Capacity string `json:"capacity"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

// FlavorAzObject contains information of the available zones
type FlavorAzObject struct {
	Capacity string   `json:"capacity"`
	AzCodes  []string `json:"az_codes"`
}

// ListResult contains the body of getting detailed
type ListResult struct {
	golangsdk.Result
}

// Extract from ListResult
func (r ListResult) Extract() ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := r.Result.ExtractInto(&s)
	return s.Flavors, err
}
