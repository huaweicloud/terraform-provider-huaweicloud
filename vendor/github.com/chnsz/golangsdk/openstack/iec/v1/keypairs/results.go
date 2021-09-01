package keypairs

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

func (r CreateResult) Extract() (*common.KeyPair, error) {
	var entity common.KeyPair
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*common.KeyPair, error) {
	var entity common.KeyPair
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*[]common.KeyPair, error) {
	var entity []common.KeyPair

	err := r.ExtractInto(&entity)
	return &entity, err
}
