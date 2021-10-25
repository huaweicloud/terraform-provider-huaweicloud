/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package pools

import (
	"github.com/chnsz/golangsdk/pagination"
)

type PoolSummaryDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	Type        string `json:"type"`
	VpcID       string `json:"vpc_id"`
	Description string `json:"description"`
}

type Pool struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	VpcID       string `json:"vpc_id"`
	Description string `json:"description"`

	Option    PoolOption    `json:"option"`
	Hosts     []IDNameEntry `json:"hosts"`
	Instances []IDNameEntry `json:"instances"`
	Bindings  []IDNameEntry `json:"bindings"`
}

type IDNameEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PoolBinding struct {
	ID             string `json:"id"`
	LoadBalancerID string `json:"loadbalancer_id"`
	WafPoolID      string `json:"waf_pool_id"`
}

type PoolPage struct {
	pagination.PageSizeBase
}

// IsEmpty checks whether a RouteTablePage struct is empty.
func (b PoolPage) IsEmpty() (bool, error) {
	arr, err := ExtractGroups(b)
	return len(arr) == 0, err
}

func ExtractGroups(r pagination.Page) ([]Pool, error) {
	var s struct {
		Total int    `json:"total"`
		Items []Pool `json:"items"`
	}
	err := (r.(PoolPage)).ExtractInto(&s)
	return s.Items, err
}

type BindELBPage struct {
	pagination.MarkerPageBase
}

// LastMarker returns the last route table ID in a ListResult
func (b BindELBPage) LastMarker() (string, error) {
	elbs, err := ExtractBindELBs(b)
	if err != nil {
		return "", err
	}
	if len(elbs) == 0 {
		return "", nil
	}
	return elbs[len(elbs)-1].ID, nil
}

// IsEmpty checks whether a RouteTablePage struct is empty.
func (b BindELBPage) IsEmpty() (bool, error) {
	elbs, err := ExtractBindELBs(b)
	return len(elbs) == 0, err
}

func ExtractBindELBs(r pagination.Page) ([]PoolBinding, error) {
	var s struct {
		Results []PoolBinding `json:"results"`
	}
	err := (r.(BindELBPage)).ExtractInto(&s)
	return s.Results, err
}
