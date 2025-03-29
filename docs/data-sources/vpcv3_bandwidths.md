---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcv3_bandwidths"
description: |-
  Use this data source to get the list of EIP bandwidths.
---

# huaweicloud_vpcv3_bandwidths

Use this data source to get the list of EIP bandwidths.

## Example Usage

```hcl
data "huaweicloud_vpcv3_bandwidths" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the bandwidth.

* `name_like` - (Optional, String) Specifies the fuzzy query name.

* `bandwidth_type` - (Optional, String) Specifies the type of the bandwidth.
  Value options:
  + **share**: Shared Bandwidth
  + **bgp**: Dynamic BGP
  + **telcom**: China Unicom
  + **sbgp**: Static BGP

* `ingress_size` - (Optional, String) Specifies the cloud access size.

* `admin_state` - (Optional, String) Specifies the status of the bandwidth.

* `billing_info` - (Optional, String) Specifies the charging information of the bandwidth.

* `tags` - (Optional, String) Specifies the tag of the bandwidth.

* `enable_bandwidth_rules` - (Optional, String) Specifies whether bandwidth groups are enabled.
  Value options: **true**, **false**.

* `rule_quota` - (Optional, Int) Specifies the rule value of the bandwidth.

* `public_border_group` - (Optional, String) Specifies the border group of the bandwidth.

* `charge_mode` - (Optional, String) Specifies the charging of the bandwidth.
  Value options: **bandwidth**, **traffic** and **95peak_plus**.

* `size` - (Optional, String) Specifies the size of the bandwidth.

* `type` - (Optional, String) Specifies the type of the bandwidth.
  Value options:
  + **WHOLE**: shared bandwidth
  + **PER**: exclusive bandwidth

* `fields` - (Optional, List) Specifies the display fields of the bandwidth.
  Value options: **id**, **name**, **tenant_id**, **size**, **ratio_95peak_plus**, **ingress_size**, **bandwidth_type**,
  **admin_state**, **billing_info**, **charge_mode**, **type**, **publicip_info**, **enable_bandwidth_rules**,
  **rule_quota**, **bandwidth_rules**, **public_border_group**, **created_at**, **updated_at**, **lock_infos**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `eip_bandwidths` - Indicates the list of bandwidths.

  The [eip_bandwidths](#eip_bandwidths_struct) structure is documented below.

<a name="eip_bandwidths_struct"></a>
The `eip_bandwidths` block supports:

* `id` - Indicates the ID of the bandwidth.

* `name` - Indicates the bBandwidth name.

* `type` - Indicates the bandwidth type.

* `size` - Indicates the bandwidth size.

* `publicip_info` - Indicates the EIP information corresponding to the bandwidth.

  The [publicip_info](#eip_bandwidths_publicip_info_struct) structure is documented below.

* `billing_info` - Indicates the Bill information.

* `project_id` - Indicates the ID of the project to which the user belongs.

* `admin_state` - Indicates the bandwidth status.

* `tags` - Indicates the EIP tags.

* `enable_bandwidth_rules` - Indicates whether bandwidth groups are enabled.

* `bandwidth_type` - Indicates the bandwidth type.

* `ingress_size` - Indicates the network access size, in Mbit/s.

* `rule_quota` - Indicates the rule value.

* `ratio_95peak_plus` - Indicates the minimum bandwidth guarantee rate of enhanced 95.

* `public_border_group` - Indicates the bandwidth AZ attribute, which indicates the center and edge.

* `bandwidth_rules` - Indicates the bandwidth rules.

  The [bandwidth_rules](#eip_bandwidths_bandwidth_rules_struct) structure is documented below.

* `created_at` - Indicates the creation time, which is a UTC time in **YYYY-MM-DDTHH:MM:SS** format.

* `updated_at` - Indicates the update time, which is a UTC time in **YYYY-MM-DDTHH:MM:SS** format.

<a name="eip_bandwidths_publicip_info_struct"></a>
The `publicip_info` block supports:

* `publicip_address` - Indicates the elastic public IPv4 or IPv6 address.

* `publicip_id` - Indicates the unique IPv4 or IPv6 address of the elastic public network.

* `publicip_type` - Indicates the EIP type.

* `publicipv6_address` - Indicates the IPv6 address.

* `ip_version` - Indicates the IP version information.

<a name="eip_bandwidths_bandwidth_rules_struct"></a>
The `bandwidth_rules` block supports:

* `id` - Indicates the bandwidth rule ID.

* `name` - Indicates the bandwidth rule name.

* `egress_size` - Indicates the maximum outbound bandwidth, in Mbit/s.

* `egress_guarented_size` - Indicates the guaranteed outbound bandwidth, in Mbit/s.

* `publicip_info` - Indicates the EIP information corresponding to the bandwidth.

  The [publicip_info](#bandwidth_rules_publicip_info_struct) structure is documented below.

<a name="bandwidth_rules_publicip_info_struct"></a>
The `publicip_info` block supports:

* `publicip_address` - Indicates the elastic public IPv4 or IPv6 address.

* `publicip_id` - Indicates the unique IPv4 or IPv6 address of the elastic public network corresponding to the bandwidth.

* `publicip_type` - Indicates the EIP type.

* `publicipv6_address` - Indicates the IPv6 address.

* `ip_version` - Indicates the IP version information.
