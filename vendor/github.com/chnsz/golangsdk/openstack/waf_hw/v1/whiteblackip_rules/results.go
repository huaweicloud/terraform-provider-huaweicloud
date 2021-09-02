/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package whiteblackip_rules

import (
	"github.com/chnsz/golangsdk"
)

type WhiteBlackIP struct {
	// WhiteBlackIP Rule ID
	Id string `json:"id"`
	// WhiteBlackIP Rule Addr
	Addr string `json:"addr"`
	// IP address type
	White int `json:"white"`
	// Policy ID
	PolicyID string `json:"policyid"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a whiteblackip rule.
func (r commonResult) Extract() (*WhiteBlackIP, error) {
	var response WhiteBlackIP
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a WhiteBlackIP rule.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a WhiteBlackIP rule.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a WhiteBlackIP rule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
