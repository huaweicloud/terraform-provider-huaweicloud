/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package pools

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts the options for creating pools.
type CreateOpts struct {
	Name        string `json:"name" required:"true"`
	Region      string `json:"region" required:"true"`
	Type        string `json:"type" required:"true"`
	VpcID       string `json:"vpc_id" required:"true"`
	Description string `json:"description,omitempty"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*PoolSummaryDetail, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PoolSummaryDetail
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, poolID string) (*Pool, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, poolID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r Pool
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

type ListPoolOpts struct {
	Name     string `q:"name"`
	VpcID    string `q:"vpc_id"`
	Detail   bool   `q:"detail"`
	Page     string `q:"page"`
	PageSize string `q:"pagesize"`
}

func List(c *golangsdk.ServiceClient, opts ListPoolOpts) (*pagination.Pager, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	page := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := PoolPage{pagination.PageSizeBase{PageResult: r}}
		return p
	})

	return &page, nil
}

type UpdatePoolOpts struct {
	Name        string      `json:"name,omitempty"`
	Option      *PoolOption `json:"option,omitempty"`
	Description *string     `json:"description,omitempty"`
}

type PoolOption struct {
	BodyLimit      int `json:"body_limit,omitempty"`
	HeaderLimit    int `json:"header_limit,omitempty"`
	ConnectTimeout int `json:"connect_timeout,omitempty"`
	SendTimeout    int `json:"send_timeout,omitempty"`
	ReadTimeout    int `json:"read_timeout,omitempty"`
}

func Update(c *golangsdk.ServiceClient, poolID string, opts UpdatePoolOpts) (*Pool, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, poolID), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r Pool
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, poolID string) (*PoolSummaryDetail, error) {
	var rst golangsdk.Result
	_, err := c.DeleteWithResponse(resourceURL(c, poolID), &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PoolSummaryDetail
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

type addLBOpts struct {
	LoadBalancerID string `json:"loadbalancer_id" required:"true"`
}

func AddELB(c *golangsdk.ServiceClient, poolID, loadBalancerID string) (*PoolBinding, error) {
	opts := addLBOpts{LoadBalancerID: loadBalancerID}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(bindingURL(c, poolID), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r PoolBinding
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func RemoveELB(c *golangsdk.ServiceClient, poolID, bindingID string) error {
	_, err := c.Delete(bindingResourceURL(c, poolID, bindingID), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}

func ListELB(c *golangsdk.ServiceClient, poolID string) pagination.Pager {
	url := bindingURL(c, poolID)

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := BindELBPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}
