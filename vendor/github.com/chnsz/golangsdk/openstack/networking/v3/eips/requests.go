package eips

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type GetResult struct {
	golangsdk.Result
}

// Get is a method by which can get the detailed information of public ip
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

type ListOpts struct {
	// Specifies the resource ID of pagination query. If the parameter
	// is left blank, only resources on the first page are queried.
	Marker string `q:"marker"`

	// Specifies the number of records returned on each page. The
	// value ranges from 0 to [2000], of which [2000] is the site difference item.
	// The specific value is determined by the site.
	Limit int `q:"limit"`

	// Specifies the query field you want.
	Fields []string `q:"fields"`

	// Sort key, valid values are id, public_ip_address, public_ipv6_address,
	// ip_version, created_at, updated_at and public_border_group.
	SortKey string `q:"sort_key"`

	// Sorting direction, valid values are asc and desc.
	SortDir string `q:"sort_dir"`

	Id []string `q:"id"`

	// Value range: 4, 6
	IPVersion int `q:"ip_version"`

	// EIP name
	Alias []string `q:"alias"`

	PublicIp   []string `q:"public_ip_address"`
	PublicIpv6 []string `q:"public_ipv6_address"`

	// Private IP address
	PrivateIp []string `q:"vnic.private_ip_address"`

	// Associated port id
	PortId []string `q:"vnic.port_id"`

	EnterpriseProjectId string `q:"enterprise_project_id"`

	// Type, valid values are EIP and DUALSTACK.
	Type []string `q:"type"`

	// Network Type, valid values are 5_telcom, 5_union, 5_bgp, 5_sbgp, 5_ipv6 and 5_graybgp
	NetworkType []string `q:"network_type"`

	// Public Pool Name
	PublicPoolName []string `q:"public_pool_name"`

	// Status, valid values are FREEZED、DOWN、ACTIVE、ERROR.
	Status []string `q:"status"`

	// Device ID
	DeviceID []string `q:"vnic.device_id"`

	// Device owner
	DeviceOwner []string `q:"vnic.device_owner"`

	VPCID []string `q:"vnic.vpc_id"`

	// Instance type the port associated with
	InstanceType []string `q:"vnic.instance_type"`

	// Instance ID the port associated with
	InstanceID []string `q:"vnic.instance_id"`

	// Associated instance type
	AssociateInstanceType []string `q:"associate_instance_type"`

	// Associated instance ID
	AssociateInstanceID []string `q:"associate_instance_id"`

	BandwidthID         []string `q:"bandwidth.id"`
	BandwidthName       []string `q:"bandwidth.name"`
	BandwidthSize       []string `q:"bandwidth.size"`
	BandwidthShareType  []string `q:"bandwidth.share_type"`
	BandwidthChargeMode []string `q:"bandwidth.charge_mode"`

	// Public border group
	PublicBorderGroup []string `q:"public_border_group"`

	AllowShareBandwidthTypeAny []string `q:"allow_share_bandwidth_type_any"`
}

// List is a method used to query the public ips with given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]PublicIp, error) {
	url := listURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := PublicIPPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractPublicIPs(pages)
}
