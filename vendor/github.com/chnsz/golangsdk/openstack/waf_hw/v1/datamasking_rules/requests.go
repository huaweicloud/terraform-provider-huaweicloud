/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package datamasking_rules

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToDataMaskingCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new datamasking rule.
type CreateOpts struct {
	Path        string `json:"url" required:"true"`
	Category    string `json:"category" required:"true"`
	Index       string `json:"index" required:"true"`
	Description string `json:"description,omitempty"`
}

// ToDataMaskingCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToDataMaskingCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new datamasking rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToDataMaskingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToDataMaskingUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a datamasking rule.
type UpdateOpts struct {
	Path        string `json:"url" required:"true"`
	Category    string `json:"category" required:"true"`
	Index       string `json:"index" required:"true"`
	Description string `json:"description,omitempty"`
}

// ToDataMaskingUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToDataMaskingUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a rule.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID, ruleID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDataMaskingUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID, ruleID), b, nil, reqOpt)
	return
}

// Get retrieves a particular datamasking rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular datamasking rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
