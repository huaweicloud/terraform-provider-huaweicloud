---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_segment_bandwidth_associate"
description: |-
  Manages a global EIP segment bandwidth association resource within HuaweiCloud.
---

# huaweicloud_global_eip_segment_bandwidth_associate

Manages a global EIP segment bandwidth association resource within HuaweiCloud.

## Example Usage

```hcl
variable "global_eip_segment_id" {}
variable "internet_bandwidth_id" {}

resource "huaweicloud_global_eip_segment_bandwidth_associate" "test" {
  global_eip_segment_id = var.global_eip_segment_id
  internet_bandwidth_id = var.internet_bandwidth_id
}
```

## Argument Reference

The following arguments are supported:

* `global_eip_segment_id` - (Required, String, NonUpdatable) Specifies the global EIP segment ID.

* `internet_bandwidth_id` - (Required, String, NonUpdatable) Specifies the internet bandwidth ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. Same with the global EIP segment ID.

* `internet_bandwidth` - The internet bandwidth information.

  The [internet_bandwidth](#internet_bandwidth_struct) structure is documented below.

<a name="internet_bandwidth_struct"></a>
The `internet_bandwidth` block supports:

* `id` - The internet bandwidth ID.

* `size` - The internet bandwidth size.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The global EIP segment bandwidth association resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_global_eip_segment_bandwidth_associate.test <id>
```
