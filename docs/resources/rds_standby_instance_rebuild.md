---
subcategory: "RDS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_standby_instance_rebuild"
description: |-
  Manages an RDS standby instance rebuild resource within HuaweiCloud.
---

# huaweicloud_rds_standby_instance_rebuild

Manages an RDS standby instance rebuild resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_standby_instance_rebuild" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.

* `workflow_id` - The task flow ID.

* `last_rebuild_time` - The last rebuilding time.

* `next_rebuild_time` - The time when the next rebuilding is allowed

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The RDS standby instance rebuild can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_standby_instance_rebuild.test <id>
```
