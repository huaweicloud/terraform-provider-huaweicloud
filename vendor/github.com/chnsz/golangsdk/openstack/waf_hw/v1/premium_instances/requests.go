/*
 Copyright (c) Huawei Technologies Co. Ltd. 2021. All rights reserved.
*/

package premium_instances

import (
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// requestOpts defines the request header and language.
var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateInstanceOpts the parameters in the creating request body
type CreateInstanceOpts struct {
	Region        string             `json:"region" required:"true"`
	ChargeMode    int                `json:"chargemode" required:"true"`
	AvailableZone string             `json:"available_zone" required:"true"`
	Arch          string             `json:"arch" required:"true"`
	NamePrefix    string             `json:"instancename" required:"true"`
	Specification string             `json:"specification" required:"true"`
	VpcId         string             `json:"vpc_id" required:"true"`
	SubnetId      string             `json:"subnet_id" required:"true"`
	SecurityGroup []string           `json:"security_group" required:"true"`
	Count         int                `json:"count" required:"true"`
	Ipv6Enable    string             `json:"ipv6_enable,omitempty"`
	VolumeType    string             `json:"volume_type,omitempty"`
	ClusterId     string             `json:"cluster_id,omitempty"`
	PoolId        string             `json:"pool_id,omitempty"`
	ResTenant     *bool              `json:"res_tenant,omitempty"`
	CpuFlavor     string             `json:"cpu_flavor,omitempty"`
	AntiAffinity  *bool              `json:"anti_affinity,omitempty"`
	Tags          []tags.ResourceTag `json:"tags,omitempty"`
}

// ListInstanceOpts the parameters in the querying request.
type ListInstanceOpts struct {
	EnterpriseProjectId string `q:"enterprise_project_id"`
	Page                int    `q:"page"`
	PageSize            int    `q:"pagesize"`
	InstanceName        string `q:"instancename"`
}

// UpdateInstanceOpts the parameters in the updating request.
type UpdateInstanceOpts struct {
	InstanceName string `json:"instancename"`
}

// CreateInstance will create a dedicated waf instance based on the values in CreateOpts.
func CreateInstance(c *golangsdk.ServiceClient, opts CreateInstanceOpts) (*CreationgRst, error) {
	return CreateWithEpsId(c, opts, "")
}

func generateEpsIdQuery(epsId string) string {
	if len(epsId) == 0 {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func CreateWithEpsId(c *golangsdk.ServiceClient, opts CreateInstanceOpts, epsId string) (*CreationgRst, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c)+generateEpsIdQuery(epsId), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r CreationgRst
		if err = rst.ExtractInto(&r); err == nil {
			return &r, nil
		}
	}
	return nil, err
}

// GetInstance get the waf instance detail.
func GetInstance(c *golangsdk.ServiceClient, id string) (*DedicatedInstance, error) {
	return GetWithEpsId(c, id, "")
}

func GetWithEpsId(c *golangsdk.ServiceClient, id string, epsId string) (*DedicatedInstance, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, id)+generateEpsIdQuery(epsId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r DedicatedInstance
		if err = rst.ExtractInto(&r); err == nil {
			return &r, nil
		}
	}
	return nil, err
}

// ListInstance query a list of waf instance base on ListInstanceOpts
func ListInstance(c *golangsdk.ServiceClient, opts ListInstanceOpts) (*DedicatedInstanceList, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r DedicatedInstanceList
		if err = rst.ExtractInto(&r); err == nil {
			return &r, nil
		}
	}
	return nil, err
}

// UpdateInstance query a list of waf instance base on UpdateInstanceOpts
func UpdateInstance(c *golangsdk.ServiceClient, id string, opts UpdateInstanceOpts) (*DedicatedInstance, error) {
	return UpdateWithEpsId(c, id, opts, "")
}

func UpdateWithEpsId(c *golangsdk.ServiceClient, id string, opts UpdateInstanceOpts, epsId string) (
	*DedicatedInstance, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, id)+generateEpsIdQuery(epsId), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r DedicatedInstance
		if err = rst.ExtractInto(&r); err == nil {
			return &r, nil
		}
	}
	return nil, err

}

// Delete will permanently delete a engine based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (*DedicatedInstance, error) {
	return DeleteWithEpsId(c, id, "")
}

func DeleteWithEpsId(c *golangsdk.ServiceClient, id string, epsId string) (*DedicatedInstance, error) {
	var rst golangsdk.Result
	_, err := c.DeleteWithResponse(resourceURL(c, id)+generateEpsIdQuery(epsId), &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders},
	)

	if err == nil {
		var r DedicatedInstance
		if err = rst.ExtractInto(&r); err == nil {
			return &r, nil
		}
	}
	return nil, err
}
