---
subcategory: "Cloud Trace Service (CTS)"
---

# huaweicloud_cts_notifications

Use this data source to get the list of CTS key event notifications.

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

* `type` - (Required, String) Specifies the type of CTS key event notifications to query. The value can be **smn** or **fun**.

* `name` - (Optional, String) Specifies the name of CTS key event notification to query.

* `status` - (Optional, String) Specifies the status of CTS key event notifications to query.
  The value can be **enabled** or **disabled**.

* `operation_type` - (Optional, String) Specifies the type of operation that will send notifications.
  The value cand be **customized** or **complete**.

* `topic_id` - (Optional, String) Specifies the URN of the topic which CTS key event notification uses.

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

* `operation_users` - The list of users.
  The [operation_users](#Notifications_OperationUsers) structure is documented below.

* `status` - The status of CTS key event notification.

* `topic_id` - The URN of the topic which CTS key event notification uses.

* `filter` - Advanced filtering conditions for the CTS key event notification.
  The [filter](#Notifications_Filter) structure is documented below.

* `created_at` - The creation time of the CTS key event notification.

<a name="Notifications_Operations"></a>
The `operations` block supports:

* `service` - The type of cloud serive.

* `resource` - The type of resource.

* `trace_names` - An array of trace names.

<a name="Notifications_OperationUsers"></a>
The `operation_users` block supports:

* `group` - The IAM user group.

* `users` - The list of IAM users.

<a name="Notifications_Filter"></a>
The `filter` block supports:

* `condition` - The relation between the rules.

* `rule` - The list of filter rules.
