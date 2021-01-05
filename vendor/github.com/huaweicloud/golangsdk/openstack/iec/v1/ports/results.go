package ports

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
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
