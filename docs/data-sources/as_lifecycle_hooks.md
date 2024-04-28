---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_lifecycle_hooks"
description: ""
---

# huaweicloud_as_lifecycle_hooks

Use this data source to get a list of AS scaling lifecycle hooks within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

data "huaweicloud_as_lifecycle_hooks" "test" {
  scaling_group_id = var.scaling_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the AS scaling group ID.

* `name` - (Optional, String) Specifies the lifecycle hook name.

* `type` - (Optional, String) Specifies the lifecycle hook type. The valid values are as follows:
  + **ADD**: The hook suspends the instance when the instance is started.
  + **REMOVE**: The hook suspends the instance when the instance is terminated.

* `default_result` - (Optional, String) Specifies the default lifecycle hook callback action. This action is
  performed when the timeout duration expires. The valid values are **ABANDON** and **CONTINUE**, defaults to **ABANDON**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `lifecycle_hooks` - All AS scaling lifecycle hooks that match the filter parameters.
  The [lifecycle_hooks](#attrblock_lifecycle_hooks) structure is documented below.

<a name="attrblock_lifecycle_hooks"></a>
The `lifecycle_hooks` block supports:

* `name` - The lifecycle hook name.

* `type` - The lifecycle hook type.

* `default_result` - The default lifecycle hook callback action.

* `timeout` - The lifecycle hook timeout duration in the unit of second.

* `notification_topic_urn` - The unique URN of the notification topic in SMN.

* `notification_topic_name` - The topic name of notification topic in SMN.

* `notification_message` - The customized notification. After a notification object is configured,
  the SMN service sends your customized notification to the object.

* `created_at` - The creation time of the lifecycle hooks.
