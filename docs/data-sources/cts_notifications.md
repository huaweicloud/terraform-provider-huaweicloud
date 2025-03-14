---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_notifications"
description: |-
  Use this data source to get the list of CTS key event notifications within HuaweiCloud.
---

# huaweicloud_cts_notifications

Use this data source to get the list of CTS key event notifications within HuaweiCloud.

## Example Usage

```hcl
variable "notification_type" {}

data "huaweicloud_cts_notifications" "test" {
  type = var.notification_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CTS key event notifications.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of CTS key event notification. The value can be **smn** or **fun**.

* `name` - (Optional, String) Specifies the name of CTS key event notification.

* `status` - (Optional, String) Specifies the status of CTS key event notification.
  The value can be **enabled** or **disabled**.

* `topic_id` - (Optional, String) Specifies the URN of the topic which CTS key event notification uses.

* `notification_id` - (Optional, String) Specifies The ID of the CTS key event notification.

* `operation_type` - (Optional, String) Specifies the type of operation that will send notifications.
  The value cand be **customized** or **complete**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `notifications` - All CTS key event notifications that match the filter parameters.
  The [notifications](#Notifications) structure is documented below.

<a name="Notifications"></a>
The `notifications` block supports:

* `id` - The ID of the CTS key event notification.

* `name` - The CTS key event notification name.

* `operation_type` - The type of operation.

* `operations` - An array of operations that will trigger notifications.
  The [operations](#Notifications_Operations) structure is documented below.

* `operation_users` - An array of users. Notifications will be sent when specified users
  perform specified operations.
  The [operation_users](#Notifications_OperationUsers) structure is documented below.

* `status` - The status of CTS key event notification.

* `topic_id` - The URN of the topic which CTS key event notification uses.

* `filter` - Advanced filtering conditions for the CTS key event notification.
  The [filter](#Notifications_Filter) structure is documented below.

* `created_at` - The creation time of the CTS key event notification.

* `agency_name` - The cloud service agency name.

<a name="Notifications_Operations"></a>
The `operations` block supports:

* `service` - The type of cloud service.

* `resource` - The type of resource.

* `trace_names` - An array of trace names.

<a name="Notifications_OperationUsers"></a>
The `operation_users` block supports:

* `group` - The IAM user group.

* `users` - An array of IAM user names in the group.

<a name="Notifications_Filter"></a>
The `filter` block supports:

* `condition` - The relation between the rules.

* `rule` - The list of filter rules.
