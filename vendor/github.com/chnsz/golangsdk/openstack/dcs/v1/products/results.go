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
	ProductID        string  `json:"product_id"`
	SpecCode         string  `json:"spec_code"`
	CacheMode        string  `json:"cache_mode"`
	Engine           string  `json:"engine"`
	EngineVersion    string  `json:"engine_versions"`
	SpecDetails      string  `json:"spec_details"`
	SpecDetails2     string  `json:"spec_details2"`
	Currency         string  `json:"currency"`
	ChargingType     string  `json:"charging_type"`
	ProductType      string  `json:"product_type"`
	ProdType         string  `json:"prod_type"`
	CPUType          string  `json:"cpu_type"`
	StorageType      string  `json:"storage_type"`
	ServiceTypeCode  string  `json:"cloud_service_type_code"`
	ResourceTypeCode string  `json:"cloud_resource_type_code"`
	ReplicaCount     int     `json:"replica_count"`
	Details          Details `json:"details"`
}

type Details struct {
	Capacity       float64 `json:"capacity"`
	MaxConnections int     `json:"max_connections"`
	MaxClients     int     `json:"max_clients"`
	MaxBandwidth   int     `json:"max_bandwidth"`
	MaxInBandwidth int     `json:"max_in_bandwidth"`
	IPCount        int     `json:"tenant_ip_count"`
	ShardingNum    int     `json:"sharding_num"`
	ProxyNum       int     `json:"proxy_num"`
	DBNumber       int     `json:"db_number"`
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
