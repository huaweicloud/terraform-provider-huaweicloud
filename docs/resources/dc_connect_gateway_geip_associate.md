---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_connect_gateway_geip_associate"
description: |-
  Manages a DC connect gateway resource within HuaweiCloud.
---

# huaweicloud_dc_connect_gateway_geip_associate

Manages a DC connect gateway resource within HuaweiCloud.

## Example Usage

```hcl
variable "connect_gateway_id" {}
variable "global_eip_id" {}

resource "huaweicloud_dc_connect_gateway_geip_associate" "test" {
  connect_gateway_id = var.connect_gateway_id
  global_eip_id      = var.global_eip_id
  type               = "IP_ADDRESS"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `connect_gateway_id` - (Required, String, NonUpdatable) Specifies the DC connect gateway ID.

* `global_eip_id` - (Required, String, NonUpdatable) Specifies the global EIP ID.

* `type` - (Optional, String, NonUpdatable) Specifies the subnet type of the global EIP. Value options: **IP_ADDRESS**,
  **IP_SEGMENT**. Defaults to **IP_ADDRESS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<connect_gateway_id>/<global_eip_id>`.

* `global_eip_segment_id` - Indicates the ID of the global EIP range.

* `status` - Indicates whether the global EIP has been bound.

* `error_message` - Indicates the cause of the failure to bind the global EIP.

* `cidr` - Indicates the global EIP and its subnet mask.

* `address_family` - Indicates the address family of the global EIP.

* `ie_vtep_ip` - Indicates the VTEP IP address of the CloudPond cluster.

* `created_time` - Indicates the time when the global EIP was bound.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DC connect gateway global EIP associate resource can be imported using the `connect_gateway_id` and `global_eip_id`,
separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dc_connect_gateway_geip_associate.test <connect_gateway_id>/<global_eip_id>
```
