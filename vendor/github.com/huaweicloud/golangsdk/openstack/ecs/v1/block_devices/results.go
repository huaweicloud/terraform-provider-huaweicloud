package block_devices

import (
	"encoding/json"

	"github.com/huaweicloud/golangsdk"
)

type VolumeAttachment struct {
	PciAddress string `json:"pciAddress"`
	Size       int    `json:"size"`
	BootIndex  int    `json:"-"`
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
