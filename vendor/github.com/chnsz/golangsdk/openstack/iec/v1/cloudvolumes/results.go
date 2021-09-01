package cloudvolumes

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

type GetResult struct {
	commonResult
}

// EdgeVolume 云磁盘
type EdgeVolume struct {
	Volume *common.Volume `json:"volume"`
}

func (r GetResult) Extract() (*EdgeVolume, error) {
	var entity EdgeVolume
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

type VolumeType struct {
	VolumeTypes []common.VolumeType `json:"volume_types"`
}

func (r GetResult) ExtractVolumeType() (*VolumeType, error) {
	var entity VolumeType
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

// Volumes 卷列表对象
type Volumes struct {
	Volumes []common.Volume `json:"volumes"`
	Count   int             `json:"count"`
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*Volumes, error) {
	var entity Volumes

	err := r.ExtractInto(&entity)
	return &entity, err
}
