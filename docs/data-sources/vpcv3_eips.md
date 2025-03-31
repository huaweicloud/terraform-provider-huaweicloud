---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcv3_eips"
description: |-
  Use this data source to get the list of EIPs.
---

# huaweicloud_vpcv3_eips

Use this data source to get the list of EIPs.

## Example Usage

```hcl
data "huaweicloud_vpcv3_eips" "eip" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, List) Specifies the type of the EIP.
  Value options:
  + **EIP**: EIP
  + **DUALSTACK**: dual-stack IPv6
  + **DUALSTACK_SUBNET**: Dual-stack subnet

* `alias` - (Optional, List) Specifies the alias of the EIP.

* `alias_like` - (Optional, String) Specifies the fuzzy search based on alias.

* `ip_version` - (Optional, List) Specifies the IP version of the EIP.
  Value options: **4**, **6**.

* `status` - (Optional, List) Specifies the status of the EIP.
  Value options: **FREEZED**, **DOWN**, **ACTIVE**, **ERROR**.

* `description` - (Optional, List) Specifies the description of the EIP.

* `public_ip_address` - (Optional, List) Specifies the public IP address of the EIP.

* `public_ip_address_like` - (Optional, String) Specifies the fuzzy search based on public IP address.

* `public_ipv6_address` - (Optional, List) Specifies the public IP v6 address  of the EIP.

* `public_ipv6_address_like` - (Optional, String) Specifies the fuzzy search based on public IP v6 address.

* `publicip_pool_name` - (Optional, List) Specifies the public IP pool name of the EIP.
  Value options: **5_telcom**, **5_union**, **5_bgp**, **5_sbgp**, **5_ipv6**, **5_graybgp** and pool name

* `fields` - (Optional, List) Specifies the display fields.
  Value options: **id**, **project_id**, **ip_version**, **type**, **public_ip_address**, **public_ipv6_address**, **status**,
  **description**, **created_at**, **updated_at**, **vnic**, **bandwidth**, **associate_instance_type**,
  **associate_instance_id**, **lock_status**, **billing_info**, **tags**, **enterprise_project_id**,
  **allow_share_bandwidth_types**, **public_border_group**, **alias**, **publicip_pool_name**, **publicip_pool_id**.

* `sort_key` - (Optional, String) Specifies the sort key.
  Value options: **id**, **public_ip_address**, **public_ipv6_address**, **ip_version**, **created_at**, **updated_at**,
  **public_border_group**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  Value options: **asc**, **desc**.

* `vnic_private_ip_address` - (Optional, List) Specifies the private IP address of the EIP.

* `vnic_private_ip_address_like` - (Optional, String) Specifies the fuzzy search based on private IP address

* `vnic_device_id` - (Optional, List) Specifies the device ID of vnic.

* `vnic_device_owner` - (Optional, List) Specifies the device owner of vnic.

* `vnic_vpc_id` - (Optional, List) Specifies the vpc ID of vnic.

* `vnic_port_id` - (Optional, List) Specifies the port ID of vnic.

* `vnic_device_owner_prefixlike` - (Optional, String) Specifies the fuzzy search based on device owner prefixlike.

* `vnic_instance_type` - (Optional, List) Specifies the instance type of vnic.

* `vnic_instance_id` - (Optional, List) Specifies the instance ID of vnic.

* `bandwidth_id` - (Optional, List) Specifies the ID of bandwidth.

* `bandwidth_name` - (Optional, List) Specifies the name of bandwidth.

* `bandwidth_name_like` - (Optional, List) Specifies the fuzzy search based on bandwidth name.

* `bandwidth_size` - (Optional, List) Specifies the size of bandwidth.

* `bandwidth_share_type` - (Optional, List) Specifies the share type of the EIP.

* `bandwidth_charge_mode` - (Optional, List) Specifies the charge mode of the EIP.

* `billing_info` - (Optional, List) Specifies the billing info of the EIP.

* `billing_mode` - (Optional, String) Specifies the billing mode of the EIP.
  Value options: **YEARLY_MONTHLY**, **PAY_PER_USE**.

* `associate_instance_type` - (Optional, List) Specifies the associate instance type of the EIP.
  Value options: **PORT**、**NATGW**、**ELB**、**VPN**、**ELBV1**

* `associate_instance_id` - (Optional, List) Specifies the associate instance ID of the EIP.

* `enterprise_project_id` - (Optional, List) Specifies the enterprise project ID of the EIP.

* `public_border_group` - (Optional, List) Specifies the public border group of the EIP.

* `allow_share_bandwidth_type_any` - (Optional, List) Specifies the shared bandwidth type of the EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `publicips` - Indicates the list of EIPs.

  The [publicips](#publicips_struct) structure is documented below.

<a name="publicips_struct"></a>
The `publicips` block supports:

* `id` - Indicates the ID of the EIP.

* `alias` - Indicates the name of the EIP.

* `status` - Indicates the status of the EIP.

* `type` - Indicates the type of the EIP.

* `enterprise_project_id` - Indicates the enterprise project ID of the EIP.

* `description` - Indicates the description of the EIP.

* `ip_version` - Indicates the IP version of the EIP.

* `publicip_pool_name` - Indicates the public pool name of the EIP.

* `project_id` - Indicates the project ID of the EIP.

* `billing_info` - Indicates the order information of the EIP.

* `vnic` - Indicates the port information when a public IP address is bound to a port instance.

  The [vnic](#publicips_vnic_struct) structure is documented below.

* `bandwidth` - Indicates the bandwidth bound to the public IP address.

  The [bandwidth](#publicips_bandwidth_struct) structure is documented below.

* `publicip_pool_id` - Indicates the ID of the network to which the public IP address belongs.

* `public_border_group` - Indicates the resources at the central site or edge site.

* `associate_instance_type` - Indicates the type of the instance bound to the public IP address.

* `associate_instance_id` - Indicates the ID of the instance bound to the public IP address.

* `lock_status` - Indicates the frozen status of the EIP.

* `public_ip_address` - Indicates the EIP or IPv6 port address.

* `public_ipv6_address` - Indicates the public IP v6 address of the EIP.

* `tags` - Indicates the tags of the EIP.

* `allow_share_bandwidth_types` - Indicates the list of shared bandwidth types that the public IP address can be added to.

* `cascade_delete_by_instance` - Indicates whether the EIP can be deleted synchronously with the instance.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

<a name="publicips_vnic_struct"></a>
The `vnic` block supports:

* `mac` - Indicates the port MAC address of the port instance.

* `instance_type` - Indicates the type.

* `vpc_id` - Indicates the VPC ID.

* `port_id` - Indicates the port ID.

* `device_owner` - Indicates the device owner.

* `port_profile` - Indicates the port profile information.

* `vtep` - Indicates the VTEP IP address.

* `vni` - Indicates the VXLAN ID.

* `instance_id` - Indicates the ID of the instance to which the port belongs.

* `private_ip_address` - Indicates the private IP address.

* `device_id` - Indicates the ID of the device to which the port belongs.

<a name="publicips_bandwidth_struct"></a>
The `bandwidth` block supports:

* `size` - Indicates the bandwidth size.

* `share_type` - Indicates the bandwidth share type.
  The value can be:
  + **PER**: exclusive bandwidth
  + **WHOLE**: shared bandwidth

* `charge_mode` - Indicates the charging mode.
  The value can be:
  + **bandwidth**: charging by bandwidth
  + **traffic**: charging by traffic
  + **95peak_plus**: charging by enhanced 95

* `name` - Indicates the bandwidth name.

* `billing_info` - Indicates the bill information.

* `id` - Indicates the bandwidth ID.
