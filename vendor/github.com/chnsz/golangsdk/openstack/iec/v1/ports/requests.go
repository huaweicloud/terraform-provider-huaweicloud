package ports

import (
	"net/http"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type CreateOpts struct {
	// Specifies the ID of the network to which the port belongs. The
	// network ID must be a real one in the network environment.
	NetworkId string `json:"network_id" required:"true"`

	DeviceOwner string `json:"device_owner,omitempty"`

	FixedIPs []FixIPEntity `json:"fixed_ips,omitempty"`
}

// FixIPEntity 私有IP的结构体
type FixIPEntity struct {
	SubnetID  string `json:"subnet_id,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

type CreateOptsBuilder interface {
	ToPortsCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToPortsCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "port")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPortsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, portId string) (r DeleteResult) {
	url := DeleteURL(client, portId)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, portId string) (r GetResult) {
	url := GetURL(client, portId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type UpdateOpts struct {
	// Specifies the port name. The value can contain no more than 255
	// characters. This parameter is left blank by default.
	Name string `json:"name,omitempty"`

	// Specifies the UUID of the security group.
	// This attribute is extended.
	SecurityGroups *[]string `json:"security_groups,omitempty"`

	// 1. Specifies a set of zero or more allowed address pairs. An
	// address pair consists of an IP address and MAC address. This attribute is extended.
	// For details, see parameter?allow_address_pair. 2. The IP address cannot be?0.0.0.0.
	// 3. Configure an independent security group for the port if a large CIDR block (subnet
	// mask less than 24) is configured for parameter?allowed_address_pairs.
	AllowedAddressPairs *[]common.AllowedAddressPair `json:"allowed_address_pairs,omitempty"`

	// Specifies a set of zero or more extra DHCP option pairs. An
	// option pair consists of an option value and name. This attribute is extended.
	ExtraDhcpOpts []common.ExtraDHCPOpt `json:"extra_dhcp_opts,omitempty"`
}

type UpdateOptsBuilder interface {
	ToPortsUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToPortsUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "port")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, portId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPortsUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateURL(client, portId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

type ListOpts struct {
	Limit          int      `q:"limit"`
	Offset         int      `q:"offset"`
	ID             string   `q:"id"`
	Name           string   `q:"name"`
	AdminStateUp   bool     `q:"admin_state_up"`
	DeviceID       string   `q:"device_id"`
	DeviceOwner    string   `q:"device_owner"`
	FixedIPs       []string `q:"fixed_ips"`
	MacAddress     string   `q:"mac_address"`
	NetworkID      string   `q:"network_id"`
	SecurityGroups string   `q:"security_groups"`
	Status         string   `q:"status"`
}

type ListPortsOptsBuilder interface {
	ToListPortsQuery() (string, error)
}

func (opts ListOpts) ToListPortsQuery() (string, error) {
	b, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListPortsOptsBuilder) (r ListResult) {
	listPortsURL := rootURL(client)
	if opts != nil {
		query, err := opts.ToListPortsQuery()
		if err != nil {
			r.Err = err
			return r
		}
		listPortsURL += query
	}

	_, r.Err = client.Get(listPortsURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
