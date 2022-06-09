package block_devices

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
)

// VolumeAttachment is the structure of the detail of disk mounting.
type VolumeAttachment struct {
	// Specifies the ECS ID in UUID format.
	ServerId string `json:"serverId"`
	// Specifies the EVS disk ID in UUID format.
	VolumeId string `json:"volumeId"`
	// Specifies the mount ID, which is the same as the EVS disk ID.
	// The value is in UUID format.
	Id string `json:"id"`
	// Specifies the drive letter of the EVS disk, which is the device name of the EVS disk.
	Device string `json:"device"`
	// Specifies the PCI address.
	PciAddress string `json:"pciAddress"`
	// Specifies the EVS disk size in GB.
	Size int `json:"size"`
	// Specifies the EVS disk boot sequence.
	// 0 indicates the system disk.
	// Non-0 indicates a data disk.
	BootIndex int `json:"bootIndex"`
	// Specifies the disk bus type.
	// Options: virtio and scsi
	BusType string `json:"bus"`
}

func (r *VolumeAttachment) UnmarshalJSON(b []byte) error {
	type tmp VolumeAttachment
	var s struct {
		tmp
		BootIndex interface{} `json:"bootIndex"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = VolumeAttachment(s.tmp)

	if s.BootIndex == nil {
		r.BootIndex = -1
	} else {
		r.BootIndex = int(s.BootIndex.(float64))
	}

	return err
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*VolumeAttachment, error) {
	s := &VolumeAttachment{}
	return s, r.ExtractInto(s)
}

func (r GetResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "volumeAttachment")
}

type ErrorResponse struct {
	// Response error.
	Error Error `json:"error"`
}

type Error struct {
	// Error code.
	Code string `json:"code"`
	// Error message.
	Message string `json:"message"`
}
