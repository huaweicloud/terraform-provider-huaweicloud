---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_warm_pool"
description: |-
  Manages an AS warm pool resource within HuaweiCloud.
---

# huaweicloud_as_warm_pool

Manages an AS warm pool resource within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

resource "huaweicloud_as_warm_pool" "test" {
  scaling_group_id        = var.scaling_group_id
  min_capacity     		  = 1
  max_capacity     		  = 1
  instance_init_wait_time = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `scaling_group_id` - (Required, String, NonUpdatable) Specifies the AS group ID.

* `min_capacity` - (Optional, Int) Specifies the minimum capacity of a warm pool.

* `max_capacity` - (Optional, Int) Specifies the maximum capacity of a warm pool.

* `instance_init_wait_time` - (Optional, Int) Specifies the instance initialization waiting time, 
  in seconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the `scaling_group_id`.

* `status` - Indicates the warm pool status. The value can be **ACTIVE**, **CLOSING**, or **CLOSED**.

## Import

The AS warm pool can be imported using `id`, e.g.

```shell
$ terraform import huaweicloud_as_warm_pool.test <id>
```
