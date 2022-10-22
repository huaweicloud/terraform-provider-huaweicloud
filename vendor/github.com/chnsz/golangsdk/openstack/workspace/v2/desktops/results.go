package desktops

import "github.com/chnsz/golangsdk/openstack/common/tags"

// RequestResp is the structure that represents the API response of desktop methods request.
type RequestResp struct {
	// Job ID.
	JobId string `json:"job_id"`
}

// RequestResp is the structure that represents the API response of Create method request.
type CreateResp struct {
	RequestResp
}

// NewVolumesResp is the structure that represents the API response of NewVolumes method request.
type NewVolumesResp struct {
	RequestResp
}

// ExpandVolumesResp is the structure that represents the API response of ExpandVolumes method request.
type ExpandVolumesResp struct {
	RequestResp
}

// UpdateResp is the structure that represents the API response of UpdateProduct method request.
type UpdateResp struct {
	// Job list.
	Jobs []Job `json:"jobs"`
}

// Job is an object to specified the operation detail of desktop.
type Job struct {
	// Desktop ID.
	DesktopId string `json:"desktop_id"`
	// Job ID.
	ID string `json:"job_id"`
}

// GetResp is the structure that represents the API response of Get method request.
type GetResp struct {
	// Desktop details.
	Desktop Desktop `json:"desktop"`
}

// Desktop is the structure that represents the desktop detail.
type Desktop struct {
	// Desktop ID.
	ID string `json:"desktop_id"`
	// Desktop name.
	Name string `json:"computer_name"`
	// IP address list of desktop.
	Addresses map[string][]AddressInfo `json:"addresses"`
	// IP address list.
	IpAddresses []string `json:"ip_addresses"`
	// Desktop type.
	Type string `json:"desktop_type"`
	// Desktop metadata.
	// + charging_mode: charging information, 1 means prePaid, 0 means postPaid.
	// + image_name: image name.
	// + metering.image_id: image ID.
	// + metering.resourcespeccode: resource specification code.
	// + metering.resourcetype: resource type.
	// + os_bit: operation system bit: 32 or 64.
	// + os_type: operation system type.
	// + desktop_os_version: operation system version.
	Metadata map[string]string `json:"metadata"`
	// Product information.
	Flavor FlavorInfo `json:"flavor"`
	// Desktop status.
	Status string `json:"status"`
	// Task status of desktop.
	// + scheduling: Creating and scheduling.
	// + block_device_mapping: In the process of being created, the disk is being prepared.
	// + networking: Creating, preparing for networking.
	// + spawning: Creating, creating internally.
	// + rebooting: rebooting.
	// + reboot_pending: Rebooting, the reboot is being issued.
	// + reboot_started: Rebooting, start internal reboot.
	// + rebooting_hard: Forced rebooting.
	// + reboot_pending_hard: A reboot is being issued during a forced reboot.
	// + reboot_started_hard: During a forced reboot, the internal reboot is started.
	// + rebuilding: rebuilding.
	// + rebuild_block_device_mapping: Rebuilding, preparing disk.
	// + rebuild_spawning: Rebuilding, rebuilding internally.
	// + migrating: Live migration is in progress.
	// + resize_prep: The specification is being adjusted and is in the preparation stage.
	// + resize_migrating: Adjusting the specification, in the migrating stage.
	// + resize_migrated: In adjusting the specification, the migration has been completed.
	// + resize_finish: Resize specification, resize is being completed.
	// + resize_reverting: In the resize specification, the resize is being rolled back.
	// + powering-off: Stopping.
	// + powering-on: Starting.
	// + deleting: Deleting.
	// + deleteFailed: Delete failed.
	TaskStatus string `json:"task_status"`
	// The time that the desktop was created.
	CreatedAt string `json:"created"`
	// Configuration of security groups
	SecurityGroups []SecurityGroup `json:"security_groups"`
	// The login status of the desktop.
	// + UNREGISTER: Indicates the state when the desktop is not registered (after the desktop is started, it will be
	//   automatically registered). The unregistered state also appears after shutdown.
	// + REGISTERED: After the desktop is registered, it is waiting for the user to connect.
	// + CONNECTED: Indicates that the user has successfully logged in and is using the desktop.
	// + DISCONNECTED: Indicates the state displayed after the session is disconnected from the desktop and the client,
	//   which may be caused by closing the client window or disconnecting the client from the desktop network.
	LoginStatus string `json:"login_status"`
	// User name.
	UserName string `json:"user_name"`
	// Product ID.
	ProductId string `json:"product_id"`
	// Configuration of root volume.
	RootVolume VolumeResp `json:"root_volume"`
	// Configuration of data volumes.
	DataVolumes []VolumeResp `json:"data_volumes"`
	// User group.
	UserGroup string `json:"user_group"`
	// Availability zone where the desktop is located.
	AvailabilityZone string `json:"availability_zone"`
	// Product information.
	Product ProductInfo `json:"product"`
	// OU name.
	OuName string `json:"ou_name"`
	// OS version.
	OsVersion string `json:"os_version"`
	// SID.
	SID string `json:"sid"`
	// Order ID.
	OrderId string `json:"order_id"`
	// The key/value pairs of the desktop.
	Tags []tags.ResourceTag `json:"tags"`
}

// AddressInfo is an object to specified the IP address details of desktop.
type AddressInfo struct {
	// IP address.
	Address string `json:"addr"`
	// IP address version.
	// + 4: IPv4
	// + 6: IPv6
	Version string `json:"version"`
	// MAC address.
	MacAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
	// IP address allocation method.
	// + fixed private IP address.
	// + floating Floating IP address.
	Type string `json:"OS-EXT-IPS:type"`
}

// FlavorInfo is an object to specified the flavor details of desktop.
type FlavorInfo struct {
	// Flavor ID.
	ID string `json:"id"`
	// Shortcut link information of relevant tags for the corresponding specifications of the desktop.
	Links []FlavorLinkInfo `json:"links"`
}

// FlavorLinkInfo is an object to specified the shortcut link information.
type FlavorLinkInfo struct {
	// Name of the shortcut link.
	Rel string `json:"rel"`
	// Address of the shortcut link.
	Hrel string `json:"hrel"`
}

// VolumeResp is an object to specified the volume details of root volume or data volume.
type VolumeResp struct {
	// Volume type.
	Type string `json:"type"`
	// Volume size.
	Size int `json:"size"`
	// The device name to which the volume is attached.
	Device string `json:"device"`
	// Unique ID of volume map.
	ID string `json:"id"`
	// volume ID.
	VolumeId string `json:"volume_id"`
	// The time that the volume was created.
	CreatedAt string `json:"create_time"`
	// Volume name.
	Name string `json:"display_name"`
}

// ProductInfo is an object to specified the product details.
type ProductInfo struct {
	// Product ID.
	ID string `json:"product_id"`
	// Flavor ID.
	FlavorId string `json:"flavor_id"`
	// Product type.
	Type string `json:"type"`
	// CPU number.
	CPU string `json:"cpu"`
	// Memory size.
	Memory string `json:"memory"`
	// Product description.
	Description string `json:"description"`
	// Charging info.
	ChargingMode string `json:"charge_mode"`
}
