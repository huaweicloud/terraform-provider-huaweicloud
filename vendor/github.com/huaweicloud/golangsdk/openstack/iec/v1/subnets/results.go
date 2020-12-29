package subnets

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/pagination"
)

// SubnetPage is the page returned by a pager when traversing over a collection
// of subnets.
type SubnetPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a SubnetPage struct is empty.
func (r SubnetPage) IsEmpty() (bool, error) {
	is, err := ExtractSubnets(r)
	return len(is) == 0, err
}

// ExtractSubnets accepts a Page struct, specifically a SubnetPage struct,
// and extracts the elements into a slice of Subnet structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSubnets(r pagination.Page) ([]common.Subnet, error) {
	var s struct {
		Subnets []common.Subnet `json:"subnets"`
	}
	err := (r.(SubnetPage)).ExtractInto(&s)
	return s.Subnets, err
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*common.Subnet, error) {
	var entity common.Subnet
	err := r.ExtractIntoStructPtr(&entity, "subnet")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*common.Subnet, error) {
	var entity common.Subnet
	err := r.ExtractIntoStructPtr(&entity, "subnet")
	return &entity, err
}

type UpdateResult struct {
	commonResult
}

type UpdateResp struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (r UpdateResult) Extract() (*UpdateResp, error) {
	var entity UpdateResp
	err := r.ExtractIntoStructPtr(&entity, "subnet")
	return &entity, err
}
