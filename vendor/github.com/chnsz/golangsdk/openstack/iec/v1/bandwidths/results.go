package bandwidths

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

//BandWidthObject 带宽的结构体，用于创建和更新请求
type BandWidthObject struct {
	BandWidth common.Bandwidth `json:"bandwidth"`
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*BandWidthObject, error) {
	var entity BandWidthObject
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*BandWidthObject, error) {
	var entity BandWidthObject
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

// BandWidths 共享带宽列表对象
type BandWidths struct {
	BandWidth []common.Bandwidth `json:"bandwidths"`
	Count     int                `json:"count"`
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*BandWidths, error) {
	var entity BandWidths
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
