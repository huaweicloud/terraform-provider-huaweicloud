---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_notifications"
description: |-
  Use this data source to get a list of AS group notifications.
---

# huaweicloud_as_notifications

Use this data source to get a list of AS group notifications.

## Example Usage

```hcl
variable "scaling_group_id" {}
variable "name" {}

data "huaweicloud_as_notifications" "test" {
  scaling_group_id = var.scaling_group_id
  topic_name       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the ID of the AS group to which the notifications belong.

* `topic_name` - (Optional, String) Specifies the topic name in SMN.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - All AS group notifications that match the filter parameters.

  The [topics](#topics_struct) structure is documented below.

<a name="topics_struct"></a>
The `topics` block supports:

* `topic_name` - The topic name in SMN.

* `topic_urn` - The unique topic URN in SMN.

* `events` - The notification scene list.
