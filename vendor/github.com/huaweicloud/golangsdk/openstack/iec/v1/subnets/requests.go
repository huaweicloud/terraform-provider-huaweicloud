package subnets

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
)

type CreateOpts struct {
	// Specifies the subnet name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name" required:"true"`

	// Specifies the network segment on which the subnet resides. The
	// value must be in CIDR format. The value must be within the CIDR block of the VPC. The
	// subnet mask cannot be greater than 28.
	Cidr string `json:"cidr" required:"true"`

	// Specifies the gateway of the subnet. The value must be a valid
	// IP address. The value must be an IP address in the subnet segment.
	GatewayIP string `json:"gateway_ip" required:"true"`

	// Specifies whether the DHCP function is enabled for the subnet.
	// The value can be true or false. If this parameter is left blank, it is set to true by
	// default.
	DhcpEnable *bool `json:"dhcp_enable,omitempty"`

	// Specifies the IP address of DNS server 1 on the subnet. The
	// value must be a valid IP address.
	PrimaryDNS string `json:"primary_dns,omitempty"`

	// Specifies the IP address of DNS server 2 on the subnet. The
	// value must be a valid IP address.
	SecondaryDNS string `json:"secondary_dns,omitempty"`

	// Specifies the DNS server address list of a subnet. This field
	// is required if you need to use more than two DNS servers. This parameter value is the
	// superset of both DNS server address 1 and DNS server address 2.
	DNSList []string `json:"dnsList,omitempty"`

	// Specifies the ID of the VPC to which the subnet belongs.
	VpcID string `json:"vpc_id" required:"true"`

	SiteID string `json:"site_id,omitempty" required:"true"`
}

type CreateOptsBuilder interface {
	ToCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "subnet")
	if err != nil {
		return nil, err
	}
	return b, nil
}

type VpcResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCreateMap()
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

func Delete(client *golangsdk.ServiceClient, subnetId string) (r DeleteResult) {
	deleteURL := DeleteURL(client, subnetId)

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

func Get(client *golangsdk.ServiceClient, subnetId string) (r GetResult) {
	getURL := GetURL(client, subnetId)

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

type UpdateOpts struct {
	// Specifies the subnet name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-).
	Name string `json:"name" required:"true"`

	// Specifies whether the DHCP function is enabled for the subnet.
	// The value can be true or false. If this parameter is left blank, it is set to true by
	// default.
	DhcpEnable *bool `json:"dhcp_enable,omitempty"`

	// Specifies the IP address of DNS server 1 on the subnet. The
	// value must be a valid IP address.
	PrimaryDNS string `json:"primary_dns,omitempty"`

	// Specifies the IP address of DNS server 2 on the subnet. The
	// value must be a valid IP address.
	SecondaryDNS string `json:"secondary_dns,omitempty"`

	// Specifies the DNS server address list of a subnet. This field
	// is required if you need to use more than two DNS servers. This parameter value is the
	// superset of both DNS server address 1 and DNS server address 2.
	DNSList []string `json:"dnsList,omitempty"`
}

type UpdateOptsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "subnet")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, subnetId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	updateURL := UpdateURL(client, subnetId)

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
