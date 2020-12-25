package vpcs

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// create request.
type CreateOptsBuilder interface {
	ToVPCCreateMap() (map[string]interface{}, error)
}
type CreateOpts struct {
	//vpc name
	Name string `json:"name" required:"true"`
	//cidr,172.16.0.0/12~172.248.255.0/24
	Cidr string `json:"cidr" required:"true"`

	//mode，SYSTEM or CUSTOMER,SYSTEM： system will design and create subnet when you need;
	// CUSTOMER: you should design and create by yourself
	Mode string `json:"mode" required:"true"`
}

// ToSecurityGroupsCreateMap converts CreateOpts structures to map[string]interface{}
func (opts CreateOpts) ToVPCCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "vpc")
	if err != nil {
		return nil, err
	}
	return b, nil
}

type UpdateOpts struct {
	// Specifies the name of the VPC. The name must be unique for a
	// tenant. The value is a string of no more than 64 characters and can contain digits,
	// letters, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`

	// Specifies the range of available subnets in the VPC. The value
	// must be in CIDR format, for example, 192.168.0.0/16. The value ranges from 10.0.0.0/8
	// to 10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to
	// 192.168.255.0/24.
	Cidr string `json:"cidr,omitempty"`
}

type UpdateOptsBuilder interface {
	ToVPCUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToVPCUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "vpc")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Get get vpc detail
func Get(client *golangsdk.ServiceClient, vpcID string) (r GetResult) {
	getURL := GetURL(client, vpcID)

	var resp *http.Response
	resp, r.Err = client.Get(getURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

// Create create vpc
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVPCCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	createURL := CreateURL(client)

	var resp *http.Response
	resp, r.Err = client.Post(createURL, b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

//Update update vpc info,name especially
func Update(client *golangsdk.ServiceClient, vpcID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVPCUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	updateURL := UpdateURL(client, vpcID)

	var resp *http.Response
	resp, r.Err = client.Put(updateURL, b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

//Delete delete the vpc
func Delete(client *golangsdk.ServiceClient, vpcID string) (r DeleteResult) {
	deleteURL := DeleteURL(client, vpcID)

	var resp *http.Response
	resp, r.Err = client.Delete(deleteURL, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}
