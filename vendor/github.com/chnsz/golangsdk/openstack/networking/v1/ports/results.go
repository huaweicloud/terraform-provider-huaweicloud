package ports

import (
	"github.com/chnsz/golangsdk/pagination"
)

// Port is an API response structure of the network VIP.
type Port struct {
	// Specifies the administrative state of the port.
	// The value can only be true, and the default value is true.
	AdminStateUp bool `json:"admin_state_up"`
	// Specifies the time when the network VIP was created.
	CreatedAt string `json:"created_at"`
	// Specifies the port ID, which uniquely identifies the port.
	ID string `json:"id"`
	// Specifies the port name.
	// The value can contain no more than 255 characters. This parameter is left blank by default.
	Name string `json:"name"`
	// Specifies the ID of the network to which the port belongs.
	// The network ID must be a real one in the network environment.
	NetworkId string `json:"network_id"`
	// Specifies the port MAC address.
	// The system automatically sets this parameter, and you are not allowed to configure the parameter value.
	MacAddress string `json:"mac_address"`
	// Specifies the port IP address.
	// A port supports only one fixed IP address that cannot be changed.
	FixedIps []FixedIp `json:"fixed_ips"`
	// Specifies the ID of the device to which the port belongs.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	DeviceId string `json:"device_id"`
	// Specifies the belonged device, which can be the DHCP server, router, load balancer, or Nova.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	DeviceOwner string `json:"device_owner"`
	// Specifies the project ID.
	TenantId string `json:"tenant_id"`
	// Specifies the port status. The status of a HANA SR-IOV VM port is always DOWN.
	// The value can be ACTIVE, BUILD, or DOWN.
	Status string `json:"status"`
	// Specifies the security group UUID (extended attribute).
	SecurityGroups []string `json:"security_groups"`
	// Specifies a set of zero or more allowed address pairs. An address pair consists of an IP address and MAC address.
	// The IP address cannot be 0.0.0.0/0.
	// Configure an independent security group for the port if a large CIDR block (subnet mask less than 24) is
	// configured for parameter AllowedAddressPairs.
	AllowedAddressPairs []AddressPair `json:"allowed_address_pairs"`
	// Specifies the extended option (extended attribute) of DHCP.
	ExtraDhcpOpts []ExtraDhcpOpt `json:"extra_dhcp_opts"`
	// Specifies the VIF details. Parameter ovs_hybrid_plug specifies whether the OVS/bridge hybrid mode is used.
	VifDetails VifDetail `json:"binding:vif_details"`
	// Specifies the custom information configured by users. This is an extended attribute.
	Profile interface{} `json:"binding:profile"`
	// Specifies the type of the bound vNIC. The value can be normal or direct.
	// Parameter normal indicates software switching.
	// Parameter direct indicates SR-IOV PCIe passthrough, which is not supported.
	VnicType string `json:"binding:vnic_type"`
	// Specifies the default private network domain name information of the primary NIC.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	DnsAssignment []DnsAssignment `json:"dns_assignment"`
	// Specifies the default private network DNS name of the primary NIC.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	DnsName string `json:"dns_name"`
	// Specifies the ID of the instance to which the port belongs, for example, RDS instance ID.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	InstanceId string `json:"instance_id"`
	// Specifies the type of the instance to which the port belongs, for example, RDS.
	// The system automatically sets this parameter, and you are not allowed to configure or change the parameter value.
	InstanceType string `json:"instance_type"`
	// Specifies whether the security option is enabled for the port.
	// If the option is not enabled, the security group and DHCP snooping do not take effect.
	PortSecurityEnabled bool `json:"port_security_enabled"`
	// Availability zone to which the port belongs.
	ZoneId string `json:"zone_id"`
}

// VifDetail is an Object specifying the VIF details.
type VifDetail struct {
	// If the value is true, indicating that it is the main network card of the virtual machine.
	PrimaryInterface bool `json:"primary_interface"`
}

// DnsAssignment is an Object specifying the private network domain information.
type DnsAssignment struct {
	// Specifies the hostname.
	Hostname string `json:"hostname"`
	// Specifies the IP address of the port.
	IpAddress string `json:"ip_address"`
	// Specifies the FQDN.
	Fqdn string `json:"fqdn"`
}

// PortPage is the page returned by a pager when traversing over a collection
// of network ports.
type PortPage struct {
	pagination.MarkerPageBase
}

// LastMarker method returns the last ID in a ports page.
func (p PortPage) LastMarker() (string, error) {
	pagePorts, err := ExtractPorts(p)
	if err != nil {
		return "", err
	}
	if len(pagePorts) == 0 {
		return "", nil
	}
	lastPort := pagePorts[len(pagePorts)-1]
	return lastPort.ID, nil
}

// IsEmpty method checks whether a PortPage struct is empty.
func (p PortPage) IsEmpty() (bool, error) {
	pagePorts, err := ExtractPorts(p)
	return len(pagePorts) == 0, err
}

// ExtractPorts accepts a Page struct, specifically a PortPage struct,
// and extracts the elements into a slice of Port structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPorts(r pagination.Page) ([]Port, error) {
	var s []Port
	err := r.(PortPage).Result.ExtractIntoSlicePtr(&s, "ports")
	return s, err
}
