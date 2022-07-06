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

// ListResp is the structure that represents the request response of List method.
type ListResp struct {
	// The engine type of the DMS products.
	Engine string `json:"engine"`
	// The supported product version types.
	Versions []string `json:"versions"`
	// The list of product details.
	Products []Product `json:"products"`
}

// Product is the structure that represents the details of the product specification.
type Product struct {
	// product type. The current product types are stand-alone and cluster.
	Type string `json:"type"`
	// Product ID.
	ProductId string `json:"product_id"`
	// The underlying resource type.
	EcsFlavorId string `json:"ecs_flavor_id"`
	// Billing type.
	BillingCode string `json:"billing_code"`
	// List of supported CPU architectures.
	ArchTypes []string `json:"arch_types"`
	// List of supported billing modes.
	//   monthly: yearly/monthly type.
	//   hourly: On-demand type.
	ChargingModes []string `json:"charging_mode"`
	// List of supported disk IO types.
	IOs []IOEntity `json:"ios"`
	// A list of features supported by the current specification instance.
	SupportFeatures []SupportFeatureEntity `json:"support_features"`
	// Properties of the current specification instance.
	Properties PropertiesEntity `json:"properties"`
}

// IOEntity is the structure that represents the disk IO type information.
type IOEntity struct {
	// Disk IO encoding.
	IoSpec string `json:"io_spec"`
	// Disk type.
	Type string `json:"type"`
	// Availability Zone.
	AvailableZones []string `json:"available_zones"`
	// Unavailable zone.
	UnavailableZones []string `json:"unavailable_zones"`
}

// SupportFeatureEntity is the structure that represents the features supported by the instance.
type SupportFeatureEntity struct {
	// function name.
	Name string `json:"name"`
	// Description of the function properties supported by the instance.
	Properties SupportFeaturePropertiesEntity `json:"properties"`
}

// SupportFeaturePropertiesEntity is the structure that represents the description of the functional properties
// supported by the instance.
type SupportFeaturePropertiesEntity struct {
	// The maximum number of tasks for the dump function.
	MaxTask string `json:"max_task"`
	// Minimum number of tasks for dump function.
	MinTask string `json:"min_task"`
	// Maximum number of nodes for dump function.
	MaxNode string `json:"max_node"`
	// Minimum number of nodes for dump function.
	MinNode string `json:"min_node"`
}

// PropertiesEntity is the structure that represents the properties of the current specification instance.
type PropertiesEntity struct {
	// Maximum number of partitions per broker.
	MaxPartitionPerBroker string `json:"max_partition_per_broker"`
	// The maximum number of brokers.
	MaxBroker string `json:"max_broker"`
	// Maximum storage per node. The unit is GB.
	MaxStoragePerNode string `json:"max_storage_per_node"`
	// Maximum number of consumers per broker.
	MaxConsumerPerBroker string `json:"max_consumer_per_broker"`
	// The minimum number of brokers.
	MinBroker string `json:"min_broker"`
	// Maximum bandwidth per broker.
	MaxBandwidthPerBroker string `json:"max_bandwidth_per_broker"`
	// Minimum storage per node. The unit is GB.
	MinStoragePerNode string `json:"min_storage_per_node"`
	// Maximum TPS per Broker.
	MaxTpsPerBroker string `json:"max_tps_per_broker"`
	// Product ID alias.
	ProductAlias string `json:"product_alias"`
}
