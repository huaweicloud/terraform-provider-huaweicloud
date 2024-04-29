---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_notification"
description: ""
---

# huaweicloud_as_notification

Manages an AS notification resource within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}
variable "topic_urn" {}
variable "events" {
  type = list(string)
}

resource "huaweicloud_as_notification" "test" {
  scaling_group_id = var.scaling_group_id
  topic_urn        = var.topic_urn
  events           = var.events
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `scaling_group_id` - (Required, String, ForceNew) Specifies the AS group ID.
  Changing this creates a new AS notification.

* `topic_urn` - (Required, String, ForceNew) Specifies the unique topic URN of the SMN.
  Changing this creates a new AS notification.

* `events` - (Required, List) Specifies the topic scene of AS group. The events include `SCALING_UP`,
  `SCALING_UP_FAIL`, `SCALING_DOWN`, `SCALING_DOWN_FAIL`, `SCALING_GROUP_ABNORMAL`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique topic URN in SMN.

* `topic_name` - The topic name in SMN.

## Import

The as notification can be imported using `scaling_group_id`, `topic_urn`, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_as_notification.test <scaling_group_id>/<topic_urn>
```
