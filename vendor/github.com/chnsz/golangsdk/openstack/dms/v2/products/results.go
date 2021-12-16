package products

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
	Details          []Detail `json:"detail"`
	Name             string   `json:"name"`
	UnavailableZones []string `json:"unavailable_zones"`
	AvailableZones   []string `json:"available_zones"`
}

// Detail for dms
type Detail struct {
	Storage          string        `json:"storage"`
	ProductID        string        `json:"product_id"`
	SpecCode         string        `json:"spec_code"`
	VMSpecification  string        `json:"vm_specification"`
	ProductInfos     []ProductInfo `json:"product_info"`
	PartitionNum     string        `json:"partition_num"`
	Bandwidth        string        `json:"bandwidth"`
	Tps              string        `json:"tps"`
	IOs              []IO          `json:"io"`
	UnavailableZones []string      `json:"unavailable_zones"`
	AvailableZones   []string      `json:"available_zones"`
	EcsFlavorId      string        `json:"ecs_flavor_id"`
	ArchType         string        `json:"arch_type"`
}

// ProductInfo for dms
type ProductInfo struct {
	Storage          string   `json:"storage"`
	NodeNum          string   `json:"node_num"`
	ProductID        string   `json:"product_id"`
	SpecCode         string   `json:"spec_code"`
	IOs              []IO     `json:"io"`
	AvailableZones   []string `json:"available_zones"`
	UnavailableZones []string `json:"unavailable_zones"`
}

type IO struct {
	IOType           string   `json:"io_type"`
	StorageSpecCode  string   `json:"storage_spec_code"`
	AvailableZones   []string `json:"available_zones"`
	UnavailableZones []string `json:"unavailable_zones"`
	VolumeType       string   `json:"volume_type"`
}
