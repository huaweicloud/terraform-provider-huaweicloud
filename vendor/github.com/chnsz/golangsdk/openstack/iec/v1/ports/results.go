package ports

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*common.Port, error) {
	var entity common.Port
	err := r.ExtractIntoStructPtr(&entity, "port")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*common.Port, error) {
	var entity common.Port
	err := r.ExtractIntoStructPtr(&entity, "port")
	return &entity, err
}

type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() (*common.Port, error) {
	var entity common.Port
	err := r.ExtractIntoStructPtr(&entity, "port")
	return &entity, err
}

// Ports 端口列表对象
type Ports struct {
	Ports []common.Port `json:"ports"`
	Count int           `json:"count"`
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*Ports, error) {
	var entity Ports
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
