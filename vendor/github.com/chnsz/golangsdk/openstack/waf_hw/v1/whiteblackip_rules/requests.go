/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package whiteblackip_rules

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToWhiteBlackIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new whiteblackip rule.
type CreateOpts struct {
	Addr  string `json:"addr" required:"true"`
	White int    `json:"white,omitempty"`
}

// ToWhiteBlackIPCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToWhiteBlackIPCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new whiteblackip rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToWhiteBlackIPCreateMap()
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
	ToWhiteBlackIPUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a whiteblackip rule.
type UpdateOpts struct {
	Addr  string `json:"addr" required:"true"`
	White *int   `json:"white" required:"true"`
}

// ToWhiteBlackIPUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToWhiteBlackIPUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a rule.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID, ruleID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToWhiteBlackIPUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID, ruleID), b, nil, reqOpt)
	return
}

// Get retrieves a particular whiteblackip rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular whiteblackip rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
