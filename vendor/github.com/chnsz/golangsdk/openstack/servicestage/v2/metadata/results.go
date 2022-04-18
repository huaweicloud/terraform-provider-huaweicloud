package metadata

// Runtime is the structure that represents the detail of the ServiceStage runtime.
type Runtime struct {
	// Type.
	Type string `json:"type_name"`
	// Display name.
	DisplayName string `json:"display_name"`
	// Default container port.
	ContainerDefaultPort int `json:"container_default_port"`
	// Type description.
	TypeDesc string `json:"type_desc"`
}

// Flavor is the structure that represents the detail of the ServiceStage flavor.
type Flavor struct {
	// Flavor ID.
	ID string `json:"flavor_id"`
	// Storage size.
	StorageSize string `json:"storage_size"`
	// CPU limit.
	NumCpu string `json:"num_cpu"`
	// Initial CPU value.
	NumCpuInit string `json:"num_cpu_init"`
	// Memory limit.
	MemorySize string `json:"memory_size"`
	// Initial memory value.
	MemorySizeInit string `json:"memory_size_init"`
	// Label.
	Label string `json:"label"`
	// Whether resource specifications are customized.
	Custom bool `json:"custom"`
}
