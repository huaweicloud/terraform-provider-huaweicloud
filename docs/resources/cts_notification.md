---
subcategory: "Cloud Trace Service (CTS)"
---

# huaweicloud_cts_notification

Manages CTS key event notification resource within HuaweiCloud.

## Example Usage

### Complete Notification

```hcl
variable "topic_urn" {}

resource "huaweicloud_cts_notification" "notify" {
  name           = "keyOperate_test"
  operation_type = "complete"
  smn_topic      = var.topic_urn
}
```

### Customized Notification

```hcl
variable "topic_urn" {}

resource "huaweicloud_cts_notification" "notify" {
  name           = "keyOperate_test"
  operation_type = "customized"
  smn_topic      = var.topic_urn

  operations {
    service     = "ECS"
    resource    = "ecs"
    trace_names = ["createServer", "deleteServer"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the CTS notification resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the notification name. The value contains a maximum of 64 characters,
  and only letters, digits, underscores(_), and Chinese characters are allowed.

* `operation_type` - (Required, String) Specifies the operation type, possible options include **complete** and
  **customized**.

* `smn_topic` - (Optional, String) Specifies the URN of a topic.

* `operations` - (Optional, List) Specifies an array of operations that will trigger notifications.
  For details, see [Supported Services and Operations](https://support.huaweicloud.com/intl/en-us/usermanual-cts/cts_03_0022.html).
  The [object](#notification_operations_object) structure is documented below.

* `operation_users` - (Optional, List) Specifies an array of users. Notifications will be sent when specified users
  perform specified operations. All users are selected by default.
  The [object](#notification_operation_users_object) structure is documented below.

* `enabled` - (Optional, Bool) Specifies whether notification is enabled, defaults to true.

<a name="notification_operations_object"></a>
The `operations` block supports:

* `service` - (Required, String) Specifies the cloud service.
  
* `resource` - (Required, String) Specifies the resource type.

* `trace_names` - (Required, List) Specifies an array of trace names.

<a name="notification_operation_users_object"></a>
The `operation_users` block supports:

* `group` - (Required, String) Specifies the IAM user group name.

* `users` - (Required, List) Specifies an array of IAM users in the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the notification name.
* `notification_id` - The notification ID in UUID format.
* `status` - The notification status, the value can be **enabled** or **disabled**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `update` - Default is 5 minute.
* `delete` - Default is 5 minute.

## Import

CTS notifications can be imported using `name`, e.g.:

```
$ terraform import huaweicloud_cts_notification.tracker your_notification
```
