/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package webtamperprotection_rules

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToWebTamperCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new web tamper protection rule.
type CreateOpts struct {
	Hostname    string `json:"hostname" required:"true"`
	Url         string `json:"url" required:"true"`
	Description string `json:"description,omitempty"`
}

// ToWebTamperCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToWebTamperCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new web tamper protection rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToWebTamperCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular web tamper protection rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular web tamper protection rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
