---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_notification"
description: |-
  Manages CTS key event notification resource within HuaweiCloud.
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

### Complete Notification and Enable Filtering

```hcl
variable "topic_urn" {}

resource "huaweicloud_cts_notification" "notify" {
  name           = "keyOperate_test"
  operation_type = "complete"
  smn_topic      = var.topic_urn
  
  filter {
    condition = "AND"
    rule      = ["code = 200","resource_name = test"]
  }
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

* `name` - (Required, String) Specifies the notification name. The value contains a maximum of `64` characters,
  and only English letters, digits, underscores(_), and Chinese characters are allowed.

* `operation_type` - (Required, String) Specifies the operation type, possible options include **complete** and
  **customized**.

* `smn_topic` - (Optional, String) Specifies the URN of a topic.

* `agency_name` - (Optional, String) Specifies the cloud service agency name. The value can only be **cts_admin_trust**.

* `operations` - (Optional, List) Specifies an array of operations that will trigger notifications.
  For details, see [Supported Services and Operations](https://support.huaweicloud.com/intl/en-us/usermanual-cts/cts_03_0022.html).
  The [operations](#CTS_Notification_Operations) structure is documented below.

* `operation_users` - (Optional, List) Specifies an array of users. Notifications will be sent when specified users
  perform specified operations. All users are selected by default.
  The [operation_users](#CTS_Notification_OperationUsers) structure is documented below.

* `enabled` - (Optional, Bool) Specifies whether notification is enabled, defaults to true.

* `filter` - (Optional, List) Specifies the filtering rules for notification.
  The [filter](#CTS_Notification_Filter) structure is documented below.

<a name="CTS_Notification_Operations"></a>
The `operations` block supports:

* `service` - (Required, String) Specifies the cloud service.
  
* `resource` - (Required, String) Specifies the resource type.

* `trace_names` - (Required, List) Specifies an array of trace names.

<a name="CTS_Notification_OperationUsers"></a>
The `operation_users` block supports:

* `group` - (Required, String) Specifies the IAM user group name.

* `users` - (Required, List) Specifies an array of IAM users in the group.

<a name="CTS_Notification_Filter"></a>
The `filter` block supports:

* `condition` - (Required, String) Specifies the relationship between multiple rules. The valid values are as follows:
  + **AND**: Effective after all filtering conditions are met.
  + **OR**: Effective when any one of the conditions is met.

* `rule` - (Required, List) Specifies an array of filtering rules. It consists of three parts,
  the first part is the **key**, the second part is the **rule**, and the third part is the **value**,
  the format is: **key != value**.
  + The **key** can be: **api_version**, **code**, **trace_rating**, **trace_type**, **resource_id** and
  **resource_name**.  
  When the key is **api_version**, the value needs to follow the regular constraint: **^ (a-zA-Z0-9_ -.) {1,64}$**.  
  When the key is **code**, the length range of value is from `1` to `256`.  
  When the key is **trace_rating**, the value can be **normal**, **warning** or **incident**.  
  When the key is **trace_type**, the value can be **ConsoleAction**, **ApiCall** or **SystemAction**.  
  When the key is **resource_id**, the length range of value is from `1` to `350`.  
  When the key is **resource_name**, the length range of value is from `1` to `256`.
  + The **rule** can be: **!=** or **=**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the notification name.

* `notification_id` - The notification ID in UUID format.

* `status` - The notification status, the value can be **enabled** or **disabled**.

* `created_at` - The creation time of the notification.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

CTS notifications can be imported using `name`, e.g.:

```bash
$ terraform import huaweicloud_cts_notification.tracker your_notification
```
