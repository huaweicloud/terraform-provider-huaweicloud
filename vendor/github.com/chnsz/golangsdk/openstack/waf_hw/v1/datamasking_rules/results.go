/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package datamasking_rules

import (
	"github.com/chnsz/golangsdk"
)

type DataMasking struct {
	// DataMasking Rule ID
	Id string `json:"id"`
	// DataMaksing Rule URL
	Path string `json:"url"`
	// Masked Field
	Category string `json:"category"`
	// Masked Subfield
	Index string `json:"index"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a datamasking rule.
func (r commonResult) Extract() (*DataMasking, error) {
	var response DataMasking
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a DataMasking rule.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a DataMasking rule.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a DataMasking rule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
