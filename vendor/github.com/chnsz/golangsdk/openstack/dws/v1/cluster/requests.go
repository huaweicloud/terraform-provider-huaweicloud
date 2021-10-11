package cluster

import (
	"log"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (
	PublicBindTypeAuto         = "auto_assign"
	PublicBindTypeNotUse       = "not_use"
	PublicBindTypeBindExisting = "bind_existing"
)

type PublicIpOpts struct {
	EipID          string `json:"eip_id,omitempty"`
	PublicBindType string `json:"public_bind_type,omitempty"`
}

type CreateOpts struct {
	Name                string             `json:"name" required:"true"`
	NumberOfNode        int                `json:"number_of_node" required:"true"`
	AvailabilityZone    string             `json:"availability_zone,omitempty"`
	SubnetID            string             `json:"subnet_id" required:"true"`
	UserPwd             string             `json:"user_pwd" required:"true"`
	SecurityGroupID     string             `json:"security_group_id" required:"true"`
	PublicIp            *PublicIpOpts      `json:"public_ip,omitempty"`
	NodeType            string             `json:"node_type" required:"true"`
	VpcID               string             `json:"vpc_id" required:"true"`
	UserName            string             `json:"user_name" required:"true"`
	Port                int                `json:"port,omitempty"`         //default：8000
	NumberOfCn          *int               `json:"number_of_cn,omitempty"` //default：2
	EnterpriseProjectId string             `json:"enterprise_project_id,omitempty"`
	Tags                []tags.ResourceTag `json:"tags,omitempty"`
}

type ListOpts struct {
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

type ResetPasswordOpts struct {
	NewPassword string `json:"new_password" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "cluster")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (*CreateClusterRst, error) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] create url:%q, body=%#v", createURL(c), b)
	var rst CreateClusterRst
	_, err = c.Post(createURL(c), b, &rst, &golangsdk.RequestOpts{MoreHeaders: RequestOpts.MoreHeaders})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, id string) (*ClusterDetail, error) {
	var rsp ClusterDetailsRst
	_, err := c.Get(resourceURL(c, id), &rsp, &golangsdk.RequestOpts{MoreHeaders: RequestOpts.MoreHeaders})
	if err == nil {
		return &rsp.Cluster, nil
	}

	return nil, err
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes: []int{202},
		JSONBody: map[string]interface{}{
			"keep_last_manual_snapshot": 0,
		},
	}

	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}

func List(c *golangsdk.ServiceClient, opts ListOpts) (*ListClustersRst, error) {
	var rst ListClustersRst
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url := listURL(c) + q.String()

	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{MoreHeaders: RequestOpts.MoreHeaders})
	if err == nil {
		return &rst, nil
	}

	return nil, err
}

func ResetPassword(c *golangsdk.ServiceClient, clusterId string, opts ResetPasswordOpts) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Post(resetPasswordURL(c, clusterId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

func Resize(c *golangsdk.ServiceClient, clusterId string, scaleOutCount int) (*golangsdk.Result, error) {
	var requestBody = map[string]interface{}{
		"scale_out": map[string]interface{}{
			"count": scaleOutCount,
		},
	}

	var r golangsdk.Result
	_, err := c.Post(resizeURL(c, clusterId), requestBody, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders})
	return &r, err
}
