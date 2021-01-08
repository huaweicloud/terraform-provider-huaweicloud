package rules

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

type CreateOpts struct {
	SecurityGroupRule *common.ReqSecurityGroupRuleEntity `json:"security_group_rule"`
}

type CreateOptsBuilder interface {
	ToSecurityGroupRuleCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToSecurityGroupRuleCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecurityGroupRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(CreateURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, securityGroupRuleID string) (r DeleteResult) {
	url := DeleteURL(client, securityGroupRuleID)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, securityGroupRuleID string) (r GetResult) {
	url := GetURL(client, securityGroupRuleID)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
