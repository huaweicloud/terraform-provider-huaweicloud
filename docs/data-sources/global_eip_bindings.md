---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_bindings"
description: |-
  Use this dataSource to get the list of global EIP and instance binding relationships for the tenant.
---

# huaweicloud_global_eip_bindings

Use this dataSource to get the list of global EIP and instance binding relationships for the tenant.

## Example Usage

```hcl
data "huaweicloud_global_eip_bindings" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fields` - (Optional, List) Specifies the display fields. Each element corresponds to one query field, for example
  **geip_id**, **geip_ip_address**, **instance_type**, **instance_id**, **vnic**, **vn_list**, **public_border_group**,
  **gcbandwidth**, **version**, **created_at**, **updated_at**, and **instance_vpc_id**.

* `geip_id` - (Optional, String) Specifies the GEIP ID.

* `geip_ip_address` - (Optional, String) Specifies the GEIP IP address.

* `public_border_group` - (Optional, String) Specifies the public border group.

* `instance_type` - (Optional, String) Specifies the bound instance type.

* `instance_id` - (Optional, String) Specifies the bound instance ID.

* `instance_vpc_id` - (Optional, String) Specifies the instance VPC ID.

* `gcbandwidth_id` - (Optional, String) Specifies the global connection bandwidth ID.

* `gcbandwidth_admin_status` - (Optional, String) Specifies the global connection bandwidth admin status.

* `gcbandwidth_size` - (Optional, String) Specifies the global connection bandwidth size.

* `gcbandwidth_sla_level` - (Optional, String) Specifies the SLA level of the global connection bandwidth.

* `gcbandwidth_dscp` - (Optional, String) Specifies the DSCP value of the global connection bandwidth.

* `vnic_private_ip_address` - (Optional, String) Specifies the port private IP address.

* `vnic_vpc_id` - (Optional, String) Specifies the port VPC ID.

* `vnic_port_id` - (Optional, String) Specifies the port ID.

* `vnic_device_id` - (Optional, String) Specifies the port device ID.

* `vnic_device_owner` - (Optional, String) Specifies the port device owner.

* `vnic_device_owner_prefixlike` - (Optional, String) Specifies the prefix-like match of port device owner.

* `vnic_instance_type` - (Optional, String) Specifies the port instance type.

* `vnic_instance_id` - (Optional, String) Specifies the port instance ID.

* `sort_key` - (Optional, String) Specifies the sort fields. Supported fields include **geip_id**, **version**,
  **public_border_group**, **geip_ip_address**, **created_at**, and **updated_at**.

* `sort_dir` - (Optional, String) Specifies the sort direction. Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `geip_bindings` - The list of GEIP binding relationships.

  The [geip_bindings](#geip_bindings_struct) structure is documented below.

<a name="geip_bindings_struct"></a>
The `geip_bindings` block supports:

* `geip_id` - The GEIP ID.

* `geip_ip_address` - The GEIP IP address.

* `public_border_group` - The public border group (center site or edge site).

* `created_at` - The creation time.

* `updated_at` - The update time.

* `instance_type` - The bound instance type.

* `instance_id` - The bound instance ID.

* `version` - The GEIP version number.

* `gcbandwidth` - The backbone bandwidth object.

  The [gcbandwidth](#gcbandwidth_struct) structure is documented below.

* `vnic` - The instance port information.

  The [vnic](#vnic_struct) structure is documented below.

* `vn_list` - The GEIP virtual nexthop information.

  The [vn_list](#vn_list_struct) structure is documented below.

<a name="gcbandwidth_struct"></a>
The `gcbandwidth` block supports:

* `id` - The backbone bandwidth UUID.

* `admin_status` - The backbone bandwidth status.

* `size` - The backbone bandwidth size.

* `short_id` - The backbone bandwidth short ID.

* `sla_level` - The network tier. Valid values are **Pt**, **Au**, **Ag** and **Cu**.

* `dscp` - The DSCP value.

<a name="vnic_struct"></a>
The `vnic` block supports:

* `private_ip_address` - The port private IP address.

* `device_id` - The port device ID.

* `device_owner` - The port device owner.

* `vpc_id` - The port VPC ID.

* `port_id` - The port UUID.

* `mac` - The port MAC address.

* `vtep` - The port VTEP address.

* `vni` - The port VNI.

* `instance_id` - The instance ID of the port.

* `instance_type` - The instance type of the port.

* `port_profile` - The port profile.

<a name="vn_list_struct"></a>
The `vn_list` block supports:

* `id` - The virtual nexthop UUID.

* `owner` - The virtual nexthop owner.

* `location` - The gateway location (POD, AZ, REGION, or GLOBAL).

* `forward_mode` - The nexthop forwarding mode. Valid values are **ACTIVE-ACTIVE** and **ACTIVE-STANDBY**.

* `cluster_id` - The gateway cluster ID.

* `hash_mode` - The load balancing strategy. Valid values are **2_TUPLE**, **3_TUPLE** and **5_TUPLE**.

* `type` - The network type. Valid values are **VLAN** and **VXLAN**.

* `vni` - The network ID used with `type`.

* `nexthops` - The nexthop list.

  The [nexthops](#nexthops_struct) structure is documented below.

* `created_at` - The creation time in UTC format.

* `updated_at` - The update time in UTC format.

<a name="nexthops_struct"></a>
The `nexthops` block supports:

* `ip_address` - The nexthop IP address.

* `mode` - The active or standby mode. Valid values are **ACTIVE** and **STANDBY**.
