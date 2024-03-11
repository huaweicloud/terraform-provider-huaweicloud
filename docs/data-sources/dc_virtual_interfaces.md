---
subcategory: "Direct Connect (DC)"
---

# huaweicloud_dc_virtual_interfaces

Use this data source to get the list of DC virtual interfaces.

## Example Usage

```hcl
variable "direct_connect_id" {}

data "huaweicloud_dc_virtual_interfaces" "test" {
  direct_connect_id = var.direct_connect_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `virtual_interface_id` - (Optional, String) Specifies the ID of the virtual interface.

* `name` - (Optional, String) Specifies the name of the virtual interface.

* `status` - (Optional, String) Specifies the status of the virtual interface.

* `direct_connect_id` - (Optional, String) Specifies the ID of the direct connection associated with the virtual interface.

* `vgw_id` - (Optional, String) Specifies the ID of the virtual gateway for the virtual interface.

* `enterprise_project_id` - (Optional, String) Indicates the ID of the enterprise project
  that the virtual interface belongs to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `virtual_interfaces` - Indicates the virtual interfaces list.
  The [virtual_interfaces](#DC_virtual_interfaces) structure is documented below.

<a name="DC_virtual_interfaces"></a>
The `virtual_interfaces` block supports:

* `id` - The ID of the virtual interface.

* `name` - The name of the virtual interface.

* `bandwidth` - The bandwidth of the virtual interface.

* `created_at` - The create time of the virtual interface.

* `description` - The description of the virtual interface.

* `direct_connect_id` - The ID of the direct connection associated with the virtual interface.

* `service_type` - The type of access gateway with the virtual interface.

* `status` - The status of the virtual interface.

* `type` - The type of the virtual interface.

* `vgw_id` - The ID of the virtual gateway for the virtual interface.

* `vlan` - The VLAN connected to the user gateway of the virtual interface.

* `enable_nqa` - Does the enable nqa functionality of virtual interface.

* `enable_bfd` - Does the enable bfd functionality of virtual interface.

* `lag_id` - The link aggregation group ID associated with vif of the virtual interface.

* `device_id` - The belong device ID of the virtual interface.

* `enterprise_project_id` - The ID of the enterprise project that the virtual interface belongs to.

* `local_gateway_v4_ip` - The cloud side gateway IPv4 interface address of the virtual interface.

* `remote_gateway_v4_ip` - The customer side gateway IPv4 interface address of the virtual interface.

* `address_family` - The address cluster type of the interface.

* `local_gateway_v6_ip` - The cloud side gateway IPv6 interface address of the virtual interface.

* `remote_gateway_v6_ip` - The customer side gateway IPv6 interface address of the virtual interface.

* `remote_ep_group` - The list of remote subnets, recording the cidrs on the tenant side.

* `route_mode` - The route mode of the virtual interface.

* `asn` - The (ASN) number for the local BGP.

* `bgp_md5` - The (MD5) password for the local BGP.

* `tags` - The key/value pairs to associate with the virtual interface.
