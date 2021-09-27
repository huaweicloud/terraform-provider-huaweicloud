package flavors

import (
	"github.com/chnsz/golangsdk"
)

// Flavor represent (virtual) hardware configurations for BMS
type Flavor struct {
	// ID is the flavor's unique ID.
	ID string `json:"id"`
	// Name is the name of the flavor.
	Name string `json:"name"`
	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs string `json:"vcpus"`
	// RAM is the amount of memory, measured in MB.
	RAM int `json:"ram"`
	// Disk is the amount of root disk, measured in GB.
	Disk string `json:"disk"`
	// IsPublic indicates whether the flavor is public.
	IsPublic bool `json:"os-flavor-access:is_public"`
	// the shortcut links of the flavor.
	Links []golangsdk.Link `json:"links"`
	// the extended fields of the BMS flavor
	OsExtraSpecs OsExtraSpecs `json:"os_extra_specs"`

	// the following are reserved attributes
	Swap       string  `json:"swap"`
	Ephemeral  int     `json:"OS-FLV-EXT-DATA:ephemeral"`
	Disabled   bool    `json:"OS-FLV-DISABLED:disabled"`
	RxTxFactor float64 `json:"rxtx_factor"`
}

// OsExtraSpecs os_extra_specs struct
type OsExtraSpecs struct {
	// the resource type corresponding to the flavor. The value is ironic.
	Type string `json:"resource_type"`
	// the CPU architecture of the BMS. The value can be: x86_64 and aarch64
	CPUArch string `json:"capabilities:cpu_arch"`
	// the type of the BMS flavor in the format of flavor abbreviation.
	// For example, if the flavor name is physical.o2.medium, the flavor type is o2m.
	FlavorType string `json:"capabilities:board_type"`
	// a flavor of the Ironic type.
	HypervisorType string `json:"capabilities:hypervisor_type"`
	// whether the BMS flavor supports EVS disks: true/false
	SupportEvs string `json:"baremetal:__support_evs"`
	// the boot source of the BMS. The value can be: LocalDisk and Volume
	BootFrom string `json:"baremetal:extBootType"`
	// the maximum number of NICs on the BMS
	NetNum string `json:"baremetal:net_num"`
	// the physical CPU specifications
	CPUDetail string `json:"baremetal:cpu_detail"`
	// the physical memory specifications
	MemoryDetail string `json:"baremetal:memory_detail"`
	// the physical disk specifications
	DiskDetail string `json:"baremetal:disk_detail"`
	// the physical NIC specifications
	NetcardDetail string `json:"baremetal:netcard_detail"`
	// Specifies the status of the BMS flavor. The value can be: normal, abandon, sellout, obt, promotion
	OperationStatus string `json:"cond:operation:status"`
	// the BMS flavor status in an AZ.
	OperationAZ string `json:"cond:operation:az"`
}

// ListResult is the response from a List operation.
type ListResult struct {
	golangsdk.Result
}

// Extract provides access to the list of flavors from the List operation.
func (r ListResult) Extract() ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := r.ExtractInto(&s)
	return s.Flavors, err
}
