---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_p2c_gateways"
description: |-
  Use the data source to get the list of VPN P2C gateways.
---

# huaweicloud_vpn_p2c_gateways

Use the data source to get the list of VPN P2C gateways.

## Example Usage

```hcl
data "huaweicloud_vpn_p2c_gateways" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `p2c_vpn_gateways` - The VPN P2C gateway List.

  The [p2c_vpn_gateways](#p2c_vpn_gateways_struct) structure is documented below.

<a name="p2c_vpn_gateways_struct"></a>
The `p2c_vpn_gateways` block supports:

* `id` - The ID of a VPN P2C gateway.

* `name` - The name of a VPN P2C gateway.

* `flavor` - The specification of a VPN P2C gateway.

* `frozen_effect` - Whether a VPN P2C gateway can be deleted after being frozen.

* `connect_subnet` - The ID of the VPC subnet used by a VPN P2C gateway.

* `eip` - The EIP information.

  The [eip](#p2c_vpn_gateways_eip_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID.

* `tags` - The tag list.

  The [tags](#p2c_vpn_gateways_tags_struct) structure is documented below.

* `availability_zone_ids` - The list of availability zone ID.

* `created_at` - The creation time.

* `current_connection_number` - The number of current client connections.

* `order_id` - The order ID.

* `updated_at` - The update time.

* `status` - The status of a VPN P2C gateway.

* `vpc_id` - The ID of the VPC to which a VPN P2C gateway connects.

* `max_connection_number` - The maximum number of concurrent client connections.

<a name="p2c_vpn_gateways_eip_struct"></a>
The `eip` block supports:

* `id` - The EIP ID.

* `type` - The EIP type.

* `ip_address` - A public IPv4 address.

* `charge_mode` - The billing mode of EIP bandwidth.

* `share_type` - The bandwidth share type.

* `ip_version` - The EIP version.

* `ip_billing_info` - The EIP order information.

* `bandwidth_id` - The bandwidth ID.

* `bandwidth_size` - The bandwidth size.

* `bandwidth_name` - The bandwidth name.

* `bandwidth_billing_info` - The bandwidth order information.

<a name="p2c_vpn_gateways_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
