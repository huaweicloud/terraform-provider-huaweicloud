---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_alarm_notification"
description: |-
  Manages a DSC alarm notification resource within HuaweiCloud.
---

# huaweicloud_dsc_alarm_notification

Manages a DSC alarm notification resource within HuaweiCloud.

## Example Usage

```hcl
variable "alarm_topic_id" {}
variable "topic_urn" {}

resource "huaweicloud_dsc_alarm_notification" "test" {
  alarm_topic_id = var.alarm_topic_id
  topic_urn      = var.topic_urn
  status         = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `alarm_topic_id` - (Required, String, ForceNew) Specifies the alarm topic ID.

  Changing this will create new resource.

* `topic_urn` - (Required, String) Specifies the unique resource identifier of an SMN topic.

* `status` - (Required, Int) Specifies the alarm notification status. Valid values are:
  + `0`: Close alarm notification.
  + `1`: Open alarm notification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `alarm_topic_id`).

## Import

DSC alarm notification resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_alarm_notification.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `status`. It is generally recommended running `terraform plan` after
importing the resource. You can then decide if changes should be applied to the resource, or the resource
definition should be updated to align with the cloud. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dsc_alarm_notification" "test" {
  ...

  lifecycle {
    ignore_changes = [
      status,
    ]
  }
}
```
