/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package whiteblackip_rules

import (
	"fmt"
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
	Name        string `json:"name,omitempty"`
	Addr        string `json:"addr,omitempty"`
	White       int    `json:"white"`
	Description string `json:"description,omitempty"`
	IPGroupID   string `json:"ip_group_id,omitempty"`
}

// ToWhiteBlackIPCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToWhiteBlackIPCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new whiteblackip rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	return CreateWithEpsId(c, opts, policyID, "")
}

func CreateWithEpsId(c *golangsdk.ServiceClient, opts CreateOptsBuilder, policyID, epsID string) (r CreateResult) {
	b, err := opts.ToWhiteBlackIPCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID)+generateEpsIdQuery(epsID), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToWhiteBlackIPUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a whiteblackip rule.
type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Addr        string `json:"addr,omitempty"`
	White       *int   `json:"white" required:"true"`
	Description string `json:"description,omitempty"`
	IPGroupID   string `json:"ip_group_id,omitempty"`
}

// ToWhiteBlackIPUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToWhiteBlackIPUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a rule.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID, ruleID string, opts UpdateOptsBuilder) (r UpdateResult) {
	return UpdateWithEpsId(c, opts, policyID, ruleID, "")
}

func UpdateWithEpsId(c *golangsdk.ServiceClient, opts UpdateOptsBuilder,
	policyID, ruleID, epsID string) (r UpdateResult) {
	b, err := opts.ToWhiteBlackIPUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID, ruleID)+generateEpsIdQuery(epsID), b, nil, reqOpt)
	return
}

// Get retrieves a particular whiteblackip rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	return GetWithEpsId(c, policyID, ruleID, "")
}

func GetWithEpsId(c *golangsdk.ServiceClient, policyID, ruleID, epsID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, policyID, ruleID)+generateEpsIdQuery(epsID), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular whiteblackip rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	return DeleteWithEpsId(c, policyID, ruleID, "")
}

func DeleteWithEpsId(c *golangsdk.ServiceClient, policyID, ruleID, epsID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID)+generateEpsIdQuery(epsID), reqOpt)
	return
}

func generateEpsIdQuery(epsID string) string {
	if len(epsID) == 0 {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsID)
}
