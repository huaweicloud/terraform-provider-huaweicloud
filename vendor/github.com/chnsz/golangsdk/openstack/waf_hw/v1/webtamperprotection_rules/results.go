/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package webtamperprotection_rules

import (
	"github.com/chnsz/golangsdk"
)

type WebTamper struct {
	Id          string `json:"id"`
	PolicyID    string `json:"policyid"`
	Hostname    string `json:"hostname"`
	Url         string `json:"url"`
	TimeStamp   int64  `json:"timestamp"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a web tamper protection rule.
func (r commonResult) Extract() (*WebTamper, error) {
	var response WebTamper
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Web Tamper Protection rule.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Web Tamper Protection rule.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Web Tamper Protection rule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
