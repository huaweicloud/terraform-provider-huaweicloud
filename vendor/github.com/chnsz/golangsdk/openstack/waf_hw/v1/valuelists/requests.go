/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package valuelists

import "github.com/chnsz/golangsdk"

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts the options to create reference table.
type CreateOpts struct {
	Name        string   `json:"name" required:"true"`
	Type        string   `json:"type" required:"true"`
	Values      []string `json:"values,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Create a reference table.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*WafValueList, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r WafValueList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

//UpdateValueListOpts the options to update reference table.
type UpdateValueListOpts struct {
	Name        string   `json:"name" required:"true"`
	Values      []string `json:"values,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// Update reference table according options and id.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateValueListOpts) (*WafValueList, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, id), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r WafValueList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Get a reference table by id.
func Get(c *golangsdk.ServiceClient, id string) (*WafValueList, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, id), &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r WafValueList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// ListValueListOpts the options to query a list of reference tables.
type ListValueListOpts struct {
	Page     int `q:"page"`
	Pagesize int `q:"pagesize"`
}

// List : Query a list of reference tables according to the options.
func List(c *golangsdk.ServiceClient, opts ListValueListOpts) (*ListValueListRst, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r ListValueListRst
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Delete reference table by id.
func Delete(c *golangsdk.ServiceClient, id string) (*WafValueList, error) {
	var rst golangsdk.Result
	_, err := c.DeleteWithResponse(resourceURL(c, id), &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r WafValueList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
