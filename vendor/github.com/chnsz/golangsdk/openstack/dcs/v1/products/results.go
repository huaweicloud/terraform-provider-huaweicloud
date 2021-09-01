package products

import (
	"github.com/chnsz/golangsdk"
)

// GetResponse response
type GetResponse struct {
	Products []Product `json:"products"`
}

// Product for dcs
type Product struct {
	ProductID        string `json:"product_id"`
	SpecCode         string `json:"spec_code"`
	CacheMode        string `json:"cache_mode"`
	Engine           string `json:"engine"`
	EngineVersion    string `json:"engine_versions"`
	SpecDetails      string `json:"spec_details"`
	SpecDetails2     string `json:"spec_details2"`
	Currency         string `json:"currency"`
	ChargingType     string `json:"charging_type"`
	ProductType      string `json:"product_type"`
	ProdType         string `json:"prod_type"`
	CPUType          string `json:"cpu_type"`
	StorageType      string `json:"storage_type"`
	ServiceTypeCode  string `json:"cloud_service_type_code"`
	ResourceTypeCode string `json:"cloud_resource_type_code"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*GetResponse, error) {
	var s GetResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}
