---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_hook_instances"
description: |-
  Use this data source to get a list of AS instances hanging information.
---

# huaweicloud_as_hook_instances

Use this data source to get a list of AS instances hanging information.

## Example Usage

```hcl
variable "scaling_group_id" {}
variable "instance_id" {}

data "huaweicloud_as_hook_instances" "test" {
  scaling_group_id = var.scaling_group_id
  instance_id      = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the ID of the AS group to which the AS instances belong.

* `instance_id` - (Optional, String) Specifies the ID of the AS instance.

* `lifecycle_hook_name` - (Optional, String) Specifies the name of the lifecycle hook.

* `lifecycle_hook_status` - (Optional, String) Specifies the status of the lifecycle hook.
  The valid values are as follows:
  + **HANGING**: Suspends the instance.
  + **CONTINUE**: Continues the instance.
  + **ABANDON**: Terminates the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_hanging_info` - All lifecycle hook information about the AS instances.

  The [instance_hanging_info](#instance_hanging_info_struct) structure is documented below.

<a name="instance_hanging_info_struct"></a>
The `instance_hanging_info` block supports:

* `scaling_group_id` - The ID of the AS group to which the AS instance belongs.

* `instance_id` - The ID of the AS instance.

* `lifecycle_hook_name` - The name of the lifecycle hook.

* `lifecycle_hook_status` - The status of the lifecycle hook.

* `lifecycle_action_key` - The lifecycle action key, which determines the lifecycle callback object.

* `default_result` - The default lifecycle hook callback operation.

* `timeout` - The timeout duration, in RFC3339 format.
