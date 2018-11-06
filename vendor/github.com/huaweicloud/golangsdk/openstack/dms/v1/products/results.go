package products

import (
	"github.com/huaweicloud/golangsdk"
)

// GetResponse response
type GetResponse struct {
	Hourly  []Parameter `json:"Hourly"`
	Monthly []Parameter `json:"Monthly"`
}

// Parameter for dms
type Parameter struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Values  []Value `json:"values"`
}

// Value for dms
type Value struct {
	Details []Detail `json:"detail"`
	Name    string   `json:"name"`
}

// Detail for dms
type Detail struct {
	Storage         string        `json:"storage"`
	ProductID       string        `json:"product_id"`
	SpecCode        string        `json:"spec_code"`
	VMSpecification string        `json:"vm_specification"`
	ProductInfos    []ProductInfo `json:"product_info"`
	PartitionNum    string        `json:"partition_num"`
	Bandwidth       string        `json:"bandwidth"`
	IOs             []IO          `json:"io"`
}

// ProductInfo for dms
type ProductInfo struct {
	Storage   string `json:"storage"`
	NodeNum   string `json:"node_num"`
	ProductID string `json:"product_id"`
	SpecCode  string `json:"spec_code"`
	IOs       []IO   `json:"io"`
}

type IO struct {
	IOType          string `json:"io_type"`
	StorageSpecCode string `json:"storage_spec_code"`
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
