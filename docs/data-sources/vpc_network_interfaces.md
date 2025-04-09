---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_interfaces"
description: |-
  Use this data-source to get a list of network interfaces.
---

# huaweicloud_vpc_network_interfaces

Use this data-source to get a list of network interfaces.

## Example Usage

```hcl
data "huaweicloud_vpc_network_interfaces" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the interface.

* `interface_id` - (Optional, List) Specifies the ID of the interface.

* `network_id` - (Optional, String) Specifies the network ID of the interface.

* `mac_address` - (Optional, String) Specifies the MAC address of the interface.

* `device_id` - (Optional, String) Specifies the device ID of the interface.

* `device_owner` - (Optional, String) Specifies the device owner of the interface.

* `status` - (Optional, String) Specifies the status of the interface.
  The value can be: **ACTIVE**, **BUILD** or **DOWN**.

* `security_groups` - (Optional, List) Specifies the security group IDs of the interface.

* `fixed_ips` - (Optional, List) Filter by fixed_ips=ip_address or fixed_ips=subnet_id.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise_project ID of the interface.
  The default value is **set all_granted_eps**.

* `enable_efi` - (Optional, Bool) Specifies whether EFI is enabled .

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ports` - The list of interfaces.

  The [ports](#ports_struct) structure is documented below.

<a name="ports_struct"></a>
The `ports` block supports:

* `fixed_ips` - The interface IP addresses.

  The [fixed_ips](#ports_fixed_ips_struct) structure is documented below.

* `security_groups` - The interface security group IDs.

* `mac_address` - The MAC address of the interface.

* `status` - The interface status.

* `allowed_address_pairs` - The IP address and MAC address pairs of the interface.

  The [allowed_address_pairs](#ports_allowed_address_pairs_struct) structure is documented below.

* `instance_id` - The ID of the instance to which the interface belongs.

* `zone_id` - The AZ that the interface belongs to.

* `binding_vnic_type` - The type of the bound vNIC.

* `binding_profile` - The user-defined settings.

* `name` - The interface name.

* `network_id` - The ID of the network that the interface belongs to.

* `device_id` - The ID of the device that the interface belongs to.

* `device_owner` - The device owner.

* `extra_dhcp_opts` - The extended DHCP option.

  The [extra_dhcp_opts](#ports_extra_dhcp_opts_struct) structure is documented below.

* `port_security_enabled` - Whether the security option is enabled for the interface.

* `ipv6_bandwidth_id` - The ID of the shared bandwidth bound to the IPv6 NIC.

* `id` - The interface ID

* `dns_assignment` - The default private network domain name information of the primary NIC.

  The [dns_assignment](#ports_dns_assignment_struct) structure is documented below.

* `dns_name` - The default private network DNS name of the primary NIC.

* `binding_vif_details` - The VIF details.

  The [binding_vif_details](#ports_binding_vif_details_struct) structure is documented below.

* `instance_type` - The type of the instance to which the interface belongs.

* `enable_efi` - Whether to enable EFI.

<a name="ports_fixed_ips_struct"></a>
The `fixed_ips` block supports:

* `ip_address` - The interface IP address.

* `subnet_id` - The ID of subnet to which the interface belongs.

<a name="ports_allowed_address_pairs_struct"></a>
The `allowed_address_pairs` block supports:

* `ip_address` - The IP address.

* `mac_address` - The MAC address.

<a name="ports_extra_dhcp_opts_struct"></a>
The `extra_dhcp_opts` block supports:

* `opt_name` - The option name.

* `opt_value` - The option value.

<a name="ports_dns_assignment_struct"></a>
The `dns_assignment` block supports:

* `hostname` - The interface host name.

* `ip_address` - The interface IP address.

* `fqdn` - The private network fully qualified domain name (FQDN) of the interface.

<a name="ports_binding_vif_details_struct"></a>
The `binding_vif_details` block supports:

* `primary_interface` - Whether this is the primary NIC.

* `port_filter` - The port filter.

* `ovs_hybrid_plug` - The ovs hybrid plug.
