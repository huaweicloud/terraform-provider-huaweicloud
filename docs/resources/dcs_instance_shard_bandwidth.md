---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_shard_bandwidth"
description: |-
  Manages a DCS instance shard bandwidth resource within HuaweiCloud.
---

# huaweicloud_dcs_instance_shard_bandwidth

Manages a DCS instance shard bandwidth resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}

resource "huaweicloud_dcs_instance_shard_bandwidth" "test"{
  instance_id = var.instance_id
  group_id    = var.group_id
  bandwidth   = 1024
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the shard.

* `bandwidth` - (Required, Int, NonUpdatable) Specifies the current bandwidth (Mbit/s).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<instance_id>/<group_id>`.

* `max_bandwidth` - Indicates the maximum bandwidth (Mbit/s).

* `assured_bandwidth` - Indicates the assured bandwidth (Mbit/s).

* `updated_at` - Indicates the update time (UTC).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 10 minutes.

## Import

The DCS instance shard bandwidth can be imported using the `instance_id` and `group_id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_dcs_instance_shard_bandwidth.test <instance_id>/<group_id>
```
