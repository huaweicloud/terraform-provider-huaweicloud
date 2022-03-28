package ports

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to build a new network port.
type CreateOpts struct {
	// Specifies the ID of the network to which the port belongs.
	NetworkId string `json:"network_id" required:"true"`
	// Specifies the port name. The value can contain no more than 255 characters.
	// This parameter is left blank by default.
	Name string `json:"name,omitempty"`
	// Specifies the administrative state of the port.
	// The value can only be true, and the default value is true.
	AdminStateUp bool `json:"admin_state_up,omitempty"`
	// Specifies the device to which the port belongs.
	// Currently, only '' and neutron:VIP_PORT are supported.
	// 'neutron:VIP_PORT' indicates the port of a virtual IP address.
	DeviceOwner string `json:"device_owner,omitempty"`
	// Specifies the port IP address.
	// A port supports only one fixed IP address that cannot be changed.
	FixedIps []FixedIp `json:"fixed_ips,omitempty"`
	// Specifies the UUID of the security group.
	SecurityGroups []string `json:"security_groups,omitempty"`
	// Specifies a set of zero or more allowed address pairs.
	AllowedAddressPairs []AddressPair `json:"allowed_address_pairs,omitempty"`
	// Specifies the extended option (extended attribute) of DHCP.
	ExtraDhcpOpts []ExtraDhcpOpt `json:"extra_dhcp_opts,omitempty"`
}

// FixedIp is an Object specifying the IP information of the port.
type FixedIp struct {
	// Specifies the subnet ID.
	// You cannot change the parameter value.
	SubnetId string `json:"subnet_id,omitempty"`
	// Specifies the port IP address.
	// You cannot change the parameter value.
	IpAddress string `json:"ip_address,omitempty"`
}

// AddressPair is an Object specifying the IP/Mac addresses pair.
type AddressPair struct {
	// Specifies the IP address.
	// The IP address cannot be 0.0.0.0/0.
	// Configure an independent security group for the port if a large CIDR block (subnet mask less than 24) is
	// configured for parameter AllowedAddressPairs.
	IpAddress string `json:"ip_address" required:"true"`
	// Specifies the MAC address.
	MacAddress string `json:"mac_address,omitempty"`
}

// ExtraDhcpOpt is an Object specifying the DHCP extended properties.
type ExtraDhcpOpt struct {
	// Specifies the DHCP option name.
	// Currently, only '51' is supported to indicate the DHCP lease time.
	OptName string `json:"opt_name,omitempty"`
	// Specifies the DHCP option value.
	// When 'OptName' is '51', the parameter format is 'Xh', indicating that the DHCP lease time is X hours.
	// The value range of 'X' is '1~30000' or '-1', '-1' means the DHCP lease time is infinite.
	OptValue string `json:"opt_value,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to build a new network port.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Port, error) {
	b, err := golangsdk.BuildRequestBody(opts, "port")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r Port
		rst.ExtractIntoStructPtr(&r, "port")
		return &r, nil
	}
	return nil, err
}

// Get is a method to obtain the network port details.
func Get(c *golangsdk.ServiceClient, portID string) (*Port, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, portID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r Port
		rst.ExtractIntoStructPtr(&r, "port")
		return &r, nil
	}
	return nil, err
}

// UpdateOpts is the structure required by the Update method to update the configuration of the specified network port.
type UpdateOpts struct {
	// Specifies the ID of the network to which the port belongs.
	Name string `json:"name"`
	// Specifies the UUID of the security group.
	SecurityGroups []string `json:"security_groups"`
	// Specifies a set of zero or more allowed address pairs.
	AllowedAddressPairs []AddressPair `json:"allowed_address_pairs"`
	// Specifies the extended option (extended attribute) of DHCP.
	ExtraDhcpOpts []ExtraDhcpOpt `json:"extra_dhcp_opts"`
}

// Update is a method to update the existing network port.
func Update(c *golangsdk.ServiceClient, portID string, opts UpdateOpts) (*Port, error) {
	b, err := golangsdk.BuildRequestBody(opts, "port")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, portID), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r Port
		rst.ExtractIntoStructPtr(&r, "port")
		return &r, nil
	}
	return nil, err
}

// Delete is a method to remove an existing network port by ID.
func Delete(c *golangsdk.ServiceClient, portID string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, portID), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToPortListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
// Filtering is achieved by passing in struct field values that map to the port
// attributes you want to see returned. Marker and Limit are used for pagination.
type ListOpts struct {
	ID                  string `q:"id"`
	Name                string `q:"name"`
	NetworkID           string `q:"network_id"`
	MACAddress          string `q:"mac_address"`
	AdminStateUp        *bool  `q:"admin_state_up"`
	DeviceID            string `q:"device_id"`
	DeviceOwner         string `q:"device_owner"`
	Status              string `q:"status"`
	Marker              string `q:"marker"`
	Limit               int    `q:"limit"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
	// fixed_ips=ip_address={ip_address}&fixed_ips=subnet_id={subnet_id}
	FixedIps []string `q:"fixed_ips"`
}

// ToPortListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPortListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// ports. It accepts a ListOpts struct, which allows you to filter the
// returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToPortListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := PortPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}
